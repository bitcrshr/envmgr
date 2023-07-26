import type { PageServerLoad } from './$types';
import { BaseClient, Issuer, TokenSet, generators } from 'openid-client';
import {
	PUBLIC_AUTH0_CLIENT_ID,
	PUBLIC_AUTH0_DOMAIN,
	PUBLIC_DASHBOARD_DOMAIN,
	PUBLIC_LOCALDEV
} from '$env/static/public';
import { PRIVATE_AUTH0_CLIENT_SECRET, PRIVATE_COOKIE_ENCRYPTION_KEY } from '$env/static/private';
import { AES } from 'crypto-js';
import { Redirect, redirect } from '@sveltejs/kit';
import { nanoid } from 'nanoid';
import { DateTime } from 'luxon';

import { JWKS } from '$lib/jwks.server';
import { jwtVerify } from 'jose';

export const load: PageServerLoad = async ({ cookies }) => {
	const auth0Issuer = await Issuer.discover(`https://${PUBLIC_AUTH0_DOMAIN}`);
	const client = new auth0Issuer.Client({
		client_id: PUBLIC_AUTH0_CLIENT_ID,
		client_secret: PRIVATE_AUTH0_CLIENT_SECRET,
		redirect_uris: [`https://${PUBLIC_DASHBOARD_DOMAIN}/oauth2/callback`],
		response_types: ['code']
	});

	const refreshToken = cookies.get('refresh_token');
	const accessToken = cookies.get('access_token') ?? '';

	const state = nanoid();
	cookies.set('oauth-state', state, {
		httpOnly: true,
		secure: !PUBLIC_LOCALDEV,
		sameSite: 'lax',
		path: '/oauth2/callback'
	});

	try {
		const { payload } = await jwtVerify(accessToken, JWKS, {
			issuer: `https://${PUBLIC_AUTH0_DOMAIN}`,
			audience: 'https://api.envmgr.dev'
		});

		if (typeof payload.exp !== 'number') {
			throw '';
		}

		const exp = DateTime.fromSeconds(payload.exp);

		if (DateTime.now() >= exp) {
			throw '';
		}
	} catch (_) {
		if (!refreshToken) {
			const authUrl = getAuthUrl(client, state);

			throw redirect(302, authUrl);
		}

		try {
			const { access_token, id_token, refresh_token } = await performTokenRefresh(
				client,
				refreshToken
			);

			access_token &&
				cookies.set('access_token', access_token, {
					httpOnly: true,
					secure: !PUBLIC_LOCALDEV,
					sameSite: 'lax',
					path: '/'
				});

			id_token &&
				cookies.set('id_token', id_token, {
					httpOnly: false,
					secure: !PUBLIC_LOCALDEV,
					sameSite: 'lax',
					path: '/'
				});

			refresh_token &&
				cookies.set('refresh_token', refresh_token, {
					httpOnly: true,
					secure: !PUBLIC_LOCALDEV,
					sameSite: 'lax',
					path: '/'
				});
		} catch (_) {
			cookies.delete('access_token');
			cookies.delete('refresh_token');
			cookies.delete('id_token');

			throw redirect(302, getAuthUrl(client, state));
		}

		throw redirect(302, '/');
	}
};

async function performTokenRefresh(client: BaseClient, refreshToken: string): Promise<TokenSet> {
	return client.refresh(refreshToken);
}

function getAuthUrl(client: BaseClient, state: string) {
	return client.authorizationUrl({
		audience: 'https://api.envmgr.dev',
		scope: 'openid offline_access profile email',
		state
	});
}
