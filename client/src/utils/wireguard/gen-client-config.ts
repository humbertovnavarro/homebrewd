interface WireguardClientConfig {
  interface: string;
  privateKey: string;
  listenPort: number;
  peer: {
    publicKey: string;
    endpoint: string;
    allowedIPs: string;
  };
}

export function generateWireguardClientConfig(config: WireguardClientConfig) {
  const { privateKey, listenPort, peer } = config;
  const clientConfig = `
[Interface]
PrivateKey = ${privateKey}
ListenPort = ${listenPort}

[Peer]
PublicKey = ${peer.publicKey}
Endpoint = ${peer.endpoint}
AllowedIPs = ${peer.allowedIPs}
  `;
  return clientConfig
}
