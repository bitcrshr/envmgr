import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import {
	PUBLIC_AUTH0_CLIENT_ID,
	PUBLIC_AUTH0_DOMAIN,
	PUBLIC_DASHBOARD_DOMAIN
} from '$env/static/public';

export const load: PageServerLoad = async ({ cookies }) => {
	const idToken = cookies.get('id_token');

	const logoutUrl = new URL(`https://${PUBLIC_AUTH0_DOMAIN}/oidc/logout`);
	logoutUrl.search = new URLSearchParams({
		id_token_hint: idToken ?? '',
		client_id: PUBLIC_AUTH0_CLIENT_ID,
		post_logout_redirect_uri: `https://${PUBLIC_DASHBOARD_DOMAIN}/`
	}).toString();

	cookies.delete('access_token');
	cookies.delete('id_token');
	cookies.delete('refresh_token');

	throw redirect(302, logoutUrl.toString());
};
