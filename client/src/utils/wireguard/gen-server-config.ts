interface WireguardServerConfig {
  privateKey: string;
  listenPort: number;
  peers: Array<{
    publicKey: string;
    allowedIps: string[];
  }>;
}

export function generateWireguardServerConfig(config: WireguardServerConfig): string {
  let configString = `[Interface]
PrivateKey = ${config.privateKey}
ListenPort = ${config.listenPort}
`;

  for (const peer of config.peers) {
    configString += `
[Peer]
PublicKey = ${peer.publicKey}
AllowedIPs = ${peer.allowedIps.join(', ')}
`;
  }

  return configString;
}