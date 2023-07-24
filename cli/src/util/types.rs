use serde::{Deserialize, Serialize};

pub type BoxedError = Box<dyn std::error::Error>;
pub type BoxedSSError = Box<dyn std::error::Error + Send + Sync>;

#[derive(Debug, Serialize, Deserialize)]
pub struct JWKS {
    pub keys: Vec<jsonwebkey::JsonWebKey>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub exp: usize,
    pub iss: String,
    pub aud: String,
}

pub type OAuth2Client = oauth2::Client<
    oauth2::StandardErrorResponse<oauth2::basic::BasicErrorResponseType>,
    oauth2::StandardTokenResponse<oauth2::EmptyExtraTokenFields, oauth2::basic::BasicTokenType>,
    oauth2::basic::BasicTokenType,
    oauth2::StandardTokenIntrospectionResponse<
        oauth2::EmptyExtraTokenFields,
        oauth2::basic::BasicTokenType,
    >,
    oauth2::StandardRevocableToken,
    oauth2::StandardErrorResponse<oauth2::RevocationErrorResponseType>,
>;
