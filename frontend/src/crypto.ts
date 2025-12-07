const IV_LENGTH = 12;

function bytesToBase64(bytes: Uint8Array) {
  let binary = "";
  const chunkSize = 0x8000;
  for (let i = 0; i < bytes.length; i += chunkSize) {
    const slice = bytes.subarray(i, i + chunkSize);
    binary += String.fromCharCode(...slice);
  }
  return btoa(binary);
}

function base64ToBytes(value: string) {
  try {
    const binary = atob(value);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i += 1) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
  } catch (err) {
    throw new Error("Invalid base64 encoding");
  }
}

async function generateKey() {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, true, [
    "encrypt",
    "decrypt",
  ]);
}

async function exportKeyBase64(key: CryptoKey) {
  const raw = await crypto.subtle.exportKey("raw", key);
  return bytesToBase64(new Uint8Array(raw));
}

export async function encryptFile(data: ArrayBuffer) {
  const key = await generateKey();
  const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH));
  const cipherBuffer = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv },
    key,
    data
  );

  return {
    ciphertext: new Uint8Array(cipherBuffer),
    keyBase64: await exportKeyBase64(key),
    ivBase64: bytesToBase64(iv),
  };
}

async function importKeyFromBase64(keyBase64: string) {
  const rawKey = base64ToBytes(keyBase64);
  if (rawKey.length !== 32) {
    throw new Error("AES-256 key must be 32 bytes (base64 of 256 bits)");
  }
  return crypto.subtle.importKey("raw", rawKey, "AES-GCM", false, ["decrypt"]);
}

export async function decryptFile(options: {
  data: ArrayBuffer;
  keyBase64: string;
  ivBase64: string;
}) {
  const { data, keyBase64, ivBase64 } = options;
  const key = await importKeyFromBase64(keyBase64.trim());
  const iv = base64ToBytes(ivBase64.trim());
  if (iv.length !== IV_LENGTH) {
    throw new Error("IV must be 12 bytes for AES-GCM");
  }

  try {
    const plainBuffer = await crypto.subtle.decrypt(
      { name: "AES-GCM", iv },
      key,
      data
    );
    return new Uint8Array(plainBuffer);
  } catch (err) {
    throw new Error("Decryption failed. Check key/IV and file.");
  }
}
