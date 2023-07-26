use std::fs;

use super::types::BoxedError;

pub fn dotenv_exists(fuzzy: bool) -> Result<bool, errors::FileIOError> {
    let dir_contents = fs::read_dir("./").map_err(|e| errors::FileIOError::FsError(e))?;

    for entry_res in dir_contents {
        let entry = entry_res.map_err(|e| errors::FileIOError::FsError(e))?;

        if entry.file_name() == ".env" {
            return Ok(true);
        }

        if let Some(s) = entry.file_name().to_str() {
            if fuzzy && s.ends_with(".env") {
                return Ok(true);
            }
        }
    }

    Ok(false)
}

pub fn create_dotenv() -> Result<(), errors::FileIOError> {
    let existing = fs::File::open(".env");
    if let Ok(_) = existing {
        return Err(errors::FileIOError::DotEnvAlreadyExists);
    }

    fs::File::create(".env").map_err(|e| errors::FileIOError::FsError(e))?;

    Ok(())
}

pub fn get_dotenv() -> Result<(), errors::FileIOError> {
    let existing = fs::File::open(".env");
    if let Err(e) = existing {
        if e.kind() == std::io::ErrorKind::NotFound {
            return Err(errors::FileIOError::NoDotEnvFound);
        }
    }

    Ok(())
}

pub mod errors {
    use thiserror::Error;

    #[derive(Debug, Error)]
    pub enum FileIOError {
        #[error("cannot create .env, it already exists")]
        DotEnvAlreadyExists,

        #[error("no .env found")]
        NoDotEnvFound,

        #[from(std::io::Error)]
        #[error("an error occurred while performing a fs operation: {0}")]
        FsError(std::io::Error),
    }
}
