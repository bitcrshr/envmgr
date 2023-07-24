use keyring::Entry;
use lazy_static::lazy_static;
use once_cell::sync::OnceCell;

static ACCESS_TOKEN_ENTRY: OnceCell<Entry> = OnceCell::new();
static REFRESH_TOKEN_ENTRY: OnceCell<Entry> = OnceCell::new();

lazy_static! {
    static ref SERVICE_NAME: String = String::from("envmgr-cli");
}

pub fn get_access_token() -> Option<String> {
    let access_token_ent = ACCESS_TOKEN_ENTRY
        .get_or_try_init(|| {
            return Entry::new(SERVICE_NAME.as_str(), "access-token");
        })
        .ok()?;

    access_token_ent.get_password().ok()
}

pub fn get_refresh_token() -> Option<String> {
    let refresh_token_ent = REFRESH_TOKEN_ENTRY
        .get_or_try_init(|| {
            return Entry::new(SERVICE_NAME.as_str(), "refresh-token");
        })
        .ok()?;

    refresh_token_ent.get_password().ok()
}

pub fn set_access_token(token: &str) -> Result<(), keyring::Error> {
    let access_token_ent = ACCESS_TOKEN_ENTRY.get_or_try_init(|| {
        return Entry::new(SERVICE_NAME.as_str(), "access-token");
    })?;

    let _ = access_token_ent.delete_password();

    access_token_ent.set_password(token)
}

pub fn set_refresh_token(token: &str) -> Result<(), keyring::Error> {
    let refresh_token_ent = REFRESH_TOKEN_ENTRY.get_or_try_init(|| {
        return Entry::new(SERVICE_NAME.as_str(), "refresh-token");
    })?;

    let _ = refresh_token_ent.delete_password();

    refresh_token_ent.set_password(token)
}
