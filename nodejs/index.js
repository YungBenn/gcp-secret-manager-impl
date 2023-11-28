import { SecretManagerServiceClient } from '@google-cloud/secret-manager';

const client = new SecretManagerServiceClient();

async function addSecretVersion() {
  const secretAccount = {
    "username": "ruben",
    "password": "ruben"
  }

  const payload = Buffer.from(JSON.stringify(secretAccount), 'utf8');

  const [version] = await client.addSecretVersion({
    parent: 'projects/1061405048387/secrets/test-secret',
    payload: {
      data: payload,
    },
  });

  console.log(`Added secret version ${version.name}`);
}

async function accessSecret() {
  const [version] = await client.accessSecretVersion({
    name: 'projects/1061405048387/secrets/test-secret/versions/2',
  });

  const payload = version.payload.data.toString('utf8');
  return payload
}

// await addSecretVersion();

const env = JSON.parse(await accessSecret())
console.log(env.username);
