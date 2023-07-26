import { createRemoteJWKSet } from 'jose';
import { PUBLIC_AUTH0_DOMAIN } from '$env/static/public';

export const JWKS = await createRemoteJWKSet(
	new URL(`https://${PUBLIC_AUTH0_DOMAIN}/.well-known/jwks.json`)
);
