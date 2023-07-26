import type { PageServerLoad } from './$types';
import { Issuer, generators } from 'openid-client';
import {
	PUBLIC_AUTH0_CLIENT_ID,
	PUBLIC_AUTH0_DOMAIN,
	PUBLIC_DASHBOARD_DOMAIN,
	PUBLIC_LOCALDEV
} from '$env/static/public';
import { PRIVATE_AUTH0_CLIENT_SECRET, PRIVATE_COOKIE_ENCRYPTION_KEY } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, url, request }) => {
	const auth0Issuer = await Issuer.discover(`https://${PUBLIC_AUTH0_DOMAIN}`);
	const client = new auth0Issuer.Client({
		client_id: PUBLIC_AUTH0_CLIENT_ID,
		client_secret: PRIVATE_AUTH0_CLIENT_SECRET,
		redirect_uris: [`https://${PUBLIC_DASHBOARD_DOMAIN}/oauth2/callback`],
		response_types: ['code']
	});

	const code = url.searchParams.get('code');
	const responseState = url.searchParams.get('state');
	const cookieState = cookies.get('oauth-state');

	if (!responseState || !code || !cookieState) {
		return;
	}

	const params = client.callbackParams(request);
	const { id_token, access_token, refresh_token } = await client.callback(
		`https://${PUBLIC_DASHBOARD_DOMAIN}/oauth2/callback`,
		params,
		{ state: cookieState }
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

	throw redirect(302, '/');
};
