use chrono::{LocalResult, TimeZone, Utc};
use jsonwebtoken::{self, Algorithm};
use lazy_static::lazy_static;
use oauth2::{
    basic::BasicClient, reqwest::http_client, AuthType, AuthUrl, ClientId, DeviceAuthorizationUrl,
    RefreshToken, Scope, StandardDeviceAuthorizationResponse, TokenResponse, TokenUrl,
};

use super::types::{BoxedError, Claims, OAuth2Client, JWKS};

lazy_static! {
    static ref CLIENT_ID: String = String::from("RyFgnv5wWWXag7IRyqjgnQz4aTOyOyDH");
    static ref AUTH_PROVIDER_DOMAIN: String = String::from(
        std::env::var("AUTH_PROVIDER_DOMAIN").unwrap_or("envmgr-dev.us.auth0.com".to_string())
    );
    static ref HTTP_CLIENT: reqwest::blocking::Client = reqwest::blocking::Client::new();
    static ref SERVICE_NAME: String = String::from("envmgr-cli");
    static ref OAUTH_CLIENT: OAuth2Client = BasicClient::new(
        ClientId::new(CLIENT_ID.to_string()),
        None,
        AuthUrl::new(format!(
            "https://{}/oauth/authorize",
            AUTH_PROVIDER_DOMAIN.to_string()
        ))
        .expect(
            format!(
                "failed to parse auth url with provider domain {}",
                AUTH_PROVIDER_DOMAIN.to_string()
            )
            .as_str()
        ),
        Some(
            TokenUrl::new(format!(
                "https://{}/oauth/token",
                AUTH_PROVIDER_DOMAIN.to_string()
            ))
            .expect(
                format!(
                    "failed to parse token url with auth provder domain {}",
                    AUTH_PROVIDER_DOMAIN.to_string()
                )
                .as_str()
            )
        ),
    )
    .set_device_authorization_url(
        DeviceAuthorizationUrl::new(format!(
            "https://{}/oauth/device/code",
            AUTH_PROVIDER_DOMAIN.to_string()
        ))
        .expect(
            format!(
                "failed to parse authorization url with auth provder domain {}",
                AUTH_PROVIDER_DOMAIN.to_string()
            )
            .as_str()
        )
    )
    .set_auth_type(AuthType::RequestBody);
}

pub fn get_token() -> Result<String, BoxedError> {
    let existing_token = super::keyring::get_access_token();

    match existing_token {
        Some(token) => {
            let exp = validate_access_token(&token)?;
            if exp >= Utc::now() {
                match super::keyring::get_refresh_token() {
                    Some(tok) => return exchange_refresh_token(&tok),
                    None => return get_new_access_token(),
                }
            }

            return Ok(token);
        }

        None => return get_new_access_token(),
    }
}

fn get_new_access_token() -> Result<String, BoxedError> {
    let details: StandardDeviceAuthorizationResponse = OAUTH_CLIENT
        .exchange_device_code()?
        .add_extra_param("audience", "https://api.envmgr.dev")
        .add_scopes([
            Scope::new(String::from("profile")),
            Scope::new(String::from("offline_access")),
        ])
        .request(http_client)?;

    println!(
        "Open this URL in your browser: {}\n\nand enter this code: {}",
        details.verification_uri().to_string(),
        details.user_code().secret().to_string()
    );

    let token_result = OAUTH_CLIENT
        .exchange_device_access_token(&details)
        .request(http_client, std::thread::sleep, None)?;

    let access_token = token_result.access_token().secret();
    let refresh_token = token_result.refresh_token();

    validate_access_token(&access_token)?;

    super::keyring::set_access_token(&access_token)?;
    if let Some(tok) = refresh_token {
        super::keyring::set_refresh_token(tok.secret())?;
    }

    Ok(access_token.to_string())
}

fn exchange_refresh_token(refresh_token: &str) -> Result<String, BoxedError> {
    let result = OAUTH_CLIENT
        .exchange_refresh_token(&RefreshToken::new(String::from(refresh_token)))
        .request(http_client)?;

    let new_access_token = result.access_token().secret();
    let new_refresh_token = result.refresh_token();

    validate_access_token(&new_access_token)?;

    super::keyring::set_access_token(new_access_token)?;

    if let Some(new_refresh_token) = new_refresh_token {
        super::keyring::set_refresh_token(new_refresh_token.secret())?;
    }

    Ok(new_access_token.to_string())
}

fn validate_access_token(
    access_token: &str,
) -> Result<chrono::DateTime<Utc>, errors::ValidationError> {
    let jwks_url = format!(
        "https://{}/.well-known/jwks.json",
        AUTH_PROVIDER_DOMAIN.to_string()
    );

    let res = HTTP_CLIENT
        .get(jwks_url)
        .send()
        .map_err(|e| errors::ValidationError::JWKSFetch(e))?;

    let jwks = res
        .json::<JWKS>()
        .map_err(|e| errors::ValidationError::JWKSParse(e))?;

    let header = jsonwebtoken::decode_header(access_token)
        .map_err(|e| errors::ValidationError::JWTHeaderParse(e))?;
    let kid = header.kid.ok_or(errors::ValidationError::JWTHeaderNoKid)?;

    let jwk = jwks
        .keys
        .into_iter()
        .find(|jwk| match &jwk.key_id {
            Some(key_id) => *key_id == kid,
            None => false,
        })
        .ok_or(errors::ValidationError::NoMatchingKID)?;

    let mut validation = jsonwebtoken::Validation::new(Algorithm::RS256);
    validation.set_audience(&["https://api.envmgr.dev"]);
    validation.set_issuer(&[
        "https://envmgr-dev.us.auth0.com/",
        "https://envmgr.us.auth0.com/",
    ]);

    let decoded =
        jsonwebtoken::decode::<Claims>(access_token, &jwk.key.to_decoding_key(), &validation)
            .map_err(|e| errors::ValidationError::JWTValidationFail(e))?;

    let exp = decoded.claims.exp;
    let exp_dt = match Utc.timestamp_opt(exp as i64, 0) {
        LocalResult::None => return Err(errors::ValidationError::JWTInvalidExp(exp)),
        LocalResult::Single(dt) => dt,
        LocalResult::Ambiguous(dt, _) => dt,
    };

    Ok(exp_dt)
}

pub mod errors {
    use thiserror::Error;

    #[allow(dead_code)]
    #[derive(Debug, Error)]
    pub enum ValidationError {
        #[error("{0} is an invalid algorithm")]
        InvalidAlgorithm(String),

        #[error("{0} is an invalid audience")]
        InvalidAudience(String),

        #[error("{0} is an invalid issuer")]
        InvalidIssuer(String),

        #[error("failed to fetch JWKS: {0}")]
        JWKSFetch(reqwest::Error),

        #[error("failed to parse JWKS: {0}")]
        JWKSParse(reqwest::Error),

        #[error("failed to parse JWT header: {0}")]
        JWTHeaderParse(jsonwebtoken::errors::Error),

        #[error("JWT did not have kid field in header")]
        JWTHeaderNoKid,

        #[error("JWK did not have a matching kid for JWT")]
        NoMatchingKID,

        #[error("JWT did not pass validation: {0}")]
        JWTValidationFail(jsonwebtoken::errors::Error),

        #[error("JWT had invalid exp claim that couldn't be parsed: {0}")]
        JWTInvalidExp(usize),
    }
}
