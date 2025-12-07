import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import { decryptFile, encryptFile } from "./crypto";
import "./App.css";

type EncryptSummary = {
  key: string;
  iv: string;
  size: number;
  fileName: string;
  ttlHours: string;
  downloadLimit: string;
};

type DecryptSummary = {
  size: number;
  fileName: string;
  downloadUrl: string;
};

function formatSize(bytes: number) {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const exponent = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1
  );
  const value = bytes / 1024 ** exponent;
  return `${value.toFixed(value >= 10 || exponent === 0 ? 0 : 1)} ${
    units[exponent]
  }`;
}

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [ttlHours, setTtlHours] = useState("24");
  const [downloadLimit, setDownloadLimit] = useState("3");
  const [status, setStatus] = useState("Awaiting file...");
  const [summary, setSummary] = useState<EncryptSummary | null>(null);
  const [error, setError] = useState("");
  const [isEncrypting, setIsEncrypting] = useState(false);

  const [cipherFile, setCipherFile] = useState<File | null>(null);
  const [keyInput, setKeyInput] = useState("");
  const [ivInput, setIvInput] = useState("");
  const [decryptStatus, setDecryptStatus] = useState(
    "Awaiting encrypted file and key..."
  );
  const [decryptError, setDecryptError] = useState("");
  const [isDecrypting, setIsDecrypting] = useState(false);
  const [decryptSummary, setDecryptSummary] = useState<DecryptSummary | null>(
    null
  );

  const downloadUrlRef = useRef<string | null>(null);

  const disableEncrypt = useMemo(
    () => !file || isEncrypting,
    [file, isEncrypting]
  );

  const disableDecrypt = useMemo(
    () => !cipherFile || !keyInput.trim() || !ivInput.trim() || isDecrypting,
    [cipherFile, ivInput, isDecrypting, keyInput]
  );

  const handleEncrypt = async (event: FormEvent) => {
    event.preventDefault();
    if (!file) {
      setError("Please choose a file to encrypt.");
      return;
    }
    setError("");
    setStatus("Encrypting locally with AES-256-GCM...");
    setIsEncrypting(true);

    try {
      const bytes = await file.arrayBuffer();
      const { ciphertext, keyBase64, ivBase64 } = await encryptFile(bytes);

      setSummary({
        key: keyBase64,
        iv: ivBase64,
        size: ciphertext.byteLength,
        fileName: file.name,
        ttlHours,
        downloadLimit,
      });
      setStatus("Encrypted locally. Keep the key safe!");
    } catch (err) {
      setError("Encryption failed. Please try again.");
      setStatus("Idle");
      console.error(err);
    } finally {
      setIsEncrypting(false);
    }
  };

  const handleDecrypt = async (event: FormEvent) => {
    event.preventDefault();
    if (!cipherFile) return;

    setDecryptError("");
    setDecryptStatus("Decrypting with provided key/IV...");
    setIsDecrypting(true);

    try {
      const cipherBytes = await cipherFile.arrayBuffer();
      const plain = await decryptFile({
        data: cipherBytes,
        keyBase64: keyInput,
        ivBase64: ivInput,
      });

      if (downloadUrlRef.current) {
        URL.revokeObjectURL(downloadUrlRef.current);
      }

      const blobUrl = URL.createObjectURL(new Blob([plain]));
      downloadUrlRef.current = blobUrl;

      setDecryptSummary({
        size: plain.byteLength,
        fileName: `decrypted-${cipherFile.name}`,
        downloadUrl: blobUrl,
      });
      setDecryptStatus("Decrypted locally. Download your file below.");
    } catch (err) {
      setDecryptError(
        err instanceof Error ? err.message : "Decryption failed. Try again."
      );
      setDecryptStatus("Idle");
      console.error(err);
    } finally {
      setIsDecrypting(false);
    }
  };

  useEffect(() => {
    return () => {
      if (downloadUrlRef.current) {
        URL.revokeObjectURL(downloadUrlRef.current);
      }
    };
  }, []);

  return (
    <div className="page">
      <header className="hero">
        <p className="eyebrow">Client-side only</p>
        <h1>Encrypt & prepare a secure share</h1>
        <p className="lede">
          Upload stays in your browser. We generate an AES-256-GCM key, encrypt
          the file, and show the key + IV for you to store or send.
        </p>
      </header>

      <main className="panels">
        <section className="panel">
          <h2>Encrypt a file</h2>
          <form className="form" onSubmit={handleEncrypt}>
            <label className="field">
              <span>File to encrypt</span>
              <input
                type="file"
                onChange={(e) => setFile(e.target.files?.[0] ?? null)}
                accept="*"
                required
              />
            </label>

            <div className="grid">
              <label className="field">
                <span>Time-to-live (hours)</span>
                <input
                  type="number"
                  min="1"
                  value={ttlHours}
                  onChange={(e) => setTtlHours(e.target.value)}
                  required
                />
              </label>
              <label className="field">
                <span>Download limit</span>
                <input
                  type="number"
                  min="1"
                  value={downloadLimit}
                  onChange={(e) => setDownloadLimit(e.target.value)}
                  required
                />
              </label>
            </div>

            <button type="submit" disabled={disableEncrypt}>
              {isEncrypting ? "Encrypting..." : "Encrypt file locally"}
            </button>
            <p className="status">{error ? error : status}</p>
          </form>

          {summary && (
            <section className="result">
              <h3>Encryption details</h3>
              <ul>
                <li>
                  <strong>File:</strong> {summary.fileName} (
                  {formatSize(summary.size)})
                </li>
                <li>
                  <strong>TTL:</strong> {summary.ttlHours} hours
                </li>
                <li>
                  <strong>Download limit:</strong> {summary.downloadLimit}
                </li>
              </ul>
              <div className="tokens">
                <div>
                  <span className="label">AES-256 key (base64)</span>
                  <code>{summary.key}</code>
                </div>
                <div>
                  <span className="label">IV (base64)</span>
                  <code>{summary.iv}</code>
                </div>
              </div>
              <p className="hint">
                Store the key and IV securely. They are not sent anywhere.
              </p>
            </section>
          )}
        </section>

        <section className="panel">
          <h2>Decrypt a file</h2>
          <form className="form" onSubmit={handleDecrypt}>
            <label className="field">
              <span>Encrypted file</span>
              <input
                type="file"
                onChange={(e) => setCipherFile(e.target.files?.[0] ?? null)}
                accept="*"
                required
              />
            </label>

            <label className="field">
              <span>AES-256 key (base64)</span>
              <input
                type="text"
                value={keyInput}
                onChange={(e) => setKeyInput(e.target.value)}
                placeholder="Paste the key from the sender"
                required
              />
            </label>

            <label className="field">
              <span>IV (base64)</span>
              <input
                type="text"
                value={ivInput}
                onChange={(e) => setIvInput(e.target.value)}
                placeholder="Paste the IV from the sender"
                required
              />
            </label>

            <button type="submit" disabled={disableDecrypt}>
              {isDecrypting ? "Decrypting..." : "Decrypt file locally"}
            </button>
            <p className="status">
              {decryptError ? decryptError : decryptStatus}
            </p>
          </form>

          {decryptSummary && (
            <section className="result">
              <h3>Decryption details</h3>
              <ul>
                <li>
                  <strong>Output:</strong> {decryptSummary.fileName} (
                  {formatSize(decryptSummary.size)})
                </li>
              </ul>
              <div className="actions">
                <a
                  className="download-link"
                  href={decryptSummary.downloadUrl}
                  download={decryptSummary.fileName}
                >
                  Download decrypted file
                </a>
              </div>
              <p className="hint">
                Nothing leaves the browser during decryption.
              </p>
            </section>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;
