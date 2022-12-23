import { exec } from 'child_process';

export async function generateWireGuardKeypair(): Promise<{ publicKey: string, privateKey: string }> {
  return new Promise((resolve, reject) => {
    exec('wg genkey', (error, stdout) => {
      if (error) {
        reject(error);
        return;
      }
      const privateKey = stdout.trim();
      exec(`echo ${privateKey} | wg pubkey`, (error, stdout, stderr) => {
        if (error) {
          console.error(stderr.trim());
          reject(error);
        } else {
          const publicKey = stdout.trim();
          resolve({ publicKey, privateKey });
        }
      });
    });
  });
}
