use lazy_static::lazy_static;
use openidconnect::{core::{CoreClient, CoreProviderMetadata, CoreAuthenticationFlow}, IssuerUrl, reqwest::http_client, ClientId, PkceCodeChallenge, CsrfToken, Nonce, RedirectUrl, AuthorizationCode};
use tiny_http::Server;
use url::Host;
use std::borrow::Cow::{Borrowed, self};

lazy_static! {
    static ref KINDE_CLIENT_ID: String = String::from("ba829a24a6b94f76860f6a2062c03af9");
    static ref KINDE_DOMAIN: String = String::from(
        std::env::var("KINDE_DOMAIN")
            .unwrap_or("dev.envmgr.dev".to_string())
    );
}

static CALLBACK_SERVER_PORT: u32 = 36864;

type SendSyncError = Box<dyn std::error::Error + Send + Sync>;

pub enum LoginError {
    FailedToStartServer(SendSyncError),
    NoCodeInResponse
}


pub fn login() -> Result<(), Box<dyn std::error::Error>> {
    println!("logging in...");

    println!("kinde domain: {}", KINDE_DOMAIN.to_string());

    let provider_metadata = CoreProviderMetadata::discover(
        &IssuerUrl::new(format!("https://{}", KINDE_DOMAIN.to_string()))?,
        http_client
    )?;

    let client = CoreClient::from_provider_metadata(
        provider_metadata, ClientId::new(KINDE_CLIENT_ID.to_string()),None 
    ).set_redirect_uri(
        RedirectUrl::new(format!("http://localhost:{}/", CALLBACK_SERVER_PORT))?
    );

    let (pkce_challenge, pkce_verifier) = PkceCodeChallenge::new_random_sha256();

    let (auth_url, csrf_token, nonce) = client
        .authorize_url(
            CoreAuthenticationFlow::AuthorizationCode,
            CsrfToken::new_random,
            Nonce::new_random
        )
        .set_pkce_challenge(pkce_challenge)
        .url();

    println!("navigate to this url in a browser: {}", auth_url);

    let server = Server::http(format!("0.0.0.0:{}", CALLBACK_SERVER_PORT)).unwrap();
    let mut code: Option<String> = None;
    let mut state: Option<String> = None;

    loop {
        let request = match server.recv() {
            Ok(request) => request,
            Err(e) => return Err(Box::new(e))
        };

        let request_url = url::Url::parse(request.url())?;
        println!("got request: {}", request_url.to_string());

        if request_url.host_str() != Some(KINDE_DOMAIN.as_str()) {
            eprintln!("invalid host: {}", request_url.host_str().unwrap_or("none"));
            continue;
        }        

        if request.method() != &tiny_http::Method::Get {
            eprintln!("invalid method: {}", request.method());
            continue;
        }

        let qps = request_url.query_pairs();

        for qp in qps {
            match qp.0 {
                Borrowed("code") => {
                    code = Some(qp.1.to_string());
                },

                Borrowed("state") => {
                    state = Some(qp.1.to_string());
                }
                _ => {}
            }
        }

        break;
    }

    drop(server);

    if code.is_none() || state.is_none() {
        eprintln!("BAD RESPONSE: NO CODE OR STATE");
        return Ok(())
    }

    let auth_code = code.unwrap();
    let csrf_state = state.unwrap();

    if *csrf_token.secret() != csrf_state {
        eprintln!("BAD RESPONSE: CSRF STATE DOES NOT MATCH");
        return Ok(())
    }

    let token_response = client.exchange_code(AuthorizationCode::new(auth_code))
        .set_pkce_verifier(pkce_verifier)
        .request(http_client)?;

    println!("token response: {:?}", token_response);

    Ok(())
}


