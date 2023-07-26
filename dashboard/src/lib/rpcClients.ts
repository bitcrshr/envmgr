import { PRIVATE_API_HOST } from '$env/static/private';
import { ChannelCredentials } from '@grpc/grpc-js';
import { EnvironmentManagerClient } from './api/environment_manager.client';
import { GrpcTransport } from '@protobuf-ts/grpc-transport';
import type { RpcInterceptor } from '@protobuf-ts/runtime-rpc';

export function getClient(accessToken: string) {
	if (!accessToken) {
		throw new Error('access token was empty');
	}

	const interceptor: RpcInterceptor = {
		interceptUnary(next, method, input, options) {
			options.meta!.Authorization = `Bearer ${accessToken}`;

			return next(method, input, options);
		}
	};

	return new EnvironmentManagerClient(
		new GrpcTransport({
			host: PRIVATE_API_HOST,
			channelCredentials: ChannelCredentials.createInsecure(),
			interceptors: [interceptor]
		})
	);
}
