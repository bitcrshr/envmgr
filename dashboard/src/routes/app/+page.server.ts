import { getClient } from '$lib/rpcClients';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request, locals }) => {
	const { accessToken } = locals;

	const client = getClient(accessToken);
};
