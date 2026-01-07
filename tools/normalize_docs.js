const fs = require('fs');
const path = require('path');
const iconv = require('iconv-lite');
const jschardet = require('jschardet');

const TARGET_EXTS = ['.md', '.txt', '.csv', '.json', '.js', '.ts', '.go', '.py', '.java'];
const IGNORE_DIRS = ['node_modules', '.git', 'dist', 'build', 'obj', 'bin'];

const walkSync = (dir, filelist = []) => {
  const files = fs.readdirSync(dir);
  files.forEach(file => {
    if (IGNORE_DIRS.includes(file)) return;
    const filepath = path.join(dir, file);
    if (fs.statSync(filepath).isDirectory()) {
      filelist = walkSync(filepath, filelist);
    } else {
      if (TARGET_EXTS.includes(path.extname(file))) {
        filelist.push(filepath);
      }
    }
  });
  return filelist;
};

const convertFile = (filepath) => {
  const buffer = fs.readFileSync(filepath);
  // jschardet may fail on very small files, default to utf8 if null
  const detected = jschardet.detect(buffer);
  
  if (!detected || !detected.encoding) return;
  
  const encoding = detected.encoding;
  if (encoding.toLowerCase() === 'utf-8' && detected.confidence > 0.9) return;

  console.log(`Converting ${filepath} (From: ${encoding})...`);
  fs.writeFileSync(filepath + '.bak', buffer);
  
  try {
    const str = iconv.decode(buffer, encoding);
    const utf8Buffer = iconv.encode(str, 'utf8');
    fs.writeFileSync(filepath, utf8Buffer);
    console.log(`  ✅ Converted to UTF-8. Backup saved.`);
  } catch (e) {
    console.error(`  ❌ Conversion failed: ${e.message}`);
  }
};

console.log("🔍 Scanning for non-UTF-8 files...");
try {
    const allFiles = walkSync('.');
    allFiles.forEach(f => convertFile(f));
    console.log("✨ Normalization complete.");
} catch(e) {
    console.error("Error during normalization:", e);
}
