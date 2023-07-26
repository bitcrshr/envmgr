import type { Handle } from '@sveltejs/kit';
import { JWKS } from '$lib/jwks.server';
import { jwtVerify } from 'jose';
import { PUBLIC_AUTH0_DOMAIN } from '$env/static/public';
import { DateTime } from 'luxon';

export const handle: Handle = async ({ event, resolve }) => {
	const accessToken = event.cookies.get('access_token');
	const idToken = event.cookies.get('id_token');

	const tokensValid =
		accessToken &&
		idToken &&
		(await tokenIsValid({ token: accessToken })) &&
		(await tokenIsValid({ token: idToken, isIdToken: true }));

	switch (event.url.pathname) {
		case '/': {
			if (tokensValid) {
				return Response.redirect(event.url.origin + '/app', 302);
			}

			return resolve(event);
		}

		case '/app': {
			if (!tokensValid) {
				return Response.redirect(event.url.origin + '/', 302);
			}

			event.locals.accessToken = accessToken;
			event.locals.idToken = idToken;

			return resolve(event);
		}

		default: {
			return resolve(event);
		}
	}
};

async function tokenIsValid({
	token,
	isIdToken
}: {
	token?: string;
	isIdToken?: boolean;
}): Promise<boolean> {
	try {
		if (!token) {
			throw new Error('token was undefined');
		}

		const opts: { issuer: string; audience?: string } = {
			issuer: `https://${PUBLIC_AUTH0_DOMAIN}/`
		};
		if (!isIdToken) {
			opts.audience = 'https://api.envmgr.dev';
		}

		const { payload } = await jwtVerify(token, JWKS, opts);

		if (!payload.exp) {
			throw new Error('jwt payload did not have exp field');
		}

		const exp = DateTime.fromSeconds(payload.exp);

		if (DateTime.now() >= exp) {
			throw new Error('token expired');
		}

		return true;
	} catch (e) {
		console.warn('error while verifying jwt:\n', e, '\n', token);

		return false;
	}
}
