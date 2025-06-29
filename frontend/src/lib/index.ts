const GO_PORT = import.meta.env.WW_GO_PORT || '8080';
const GO_INTERFACE = import.meta.env.WW_GO_INTERFACE || 'localhost';

export const MAX_RECONNECT_ATTEMPTS = 5;
export const RECONNECT_DELAY = 1000;

export const GOLANG_HTTP_URL =
    import.meta.env.WW_GOLANG_HTTP_URL || `http://${GO_INTERFACE}:${GO_PORT}`;
export const GOLANG_WS_URL = import.meta.env.WW_GOLANG_WS_URL || `ws://${GO_INTERFACE}:${GO_PORT}`;
