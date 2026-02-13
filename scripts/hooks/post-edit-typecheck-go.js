#!/usr/bin/env node
/**
 * PostToolUse Hook: Go vet check after editing .go files
 *
 * Cross-platform (Windows, macOS, Linux)
 *
 * Runs after Edit tool use on Go files. Walks up from the file's
 * directory to find the nearest go.mod, then runs go vet
 * and reports only errors related to the edited file.
 */

const { execFileSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const MAX_STDIN = 1024 * 1024; // 1MB limit
let data = '';
process.stdin.setEncoding('utf8');

process.stdin.on('data', chunk => {
  if (data.length < MAX_STDIN) {
    data += chunk;
  }
});

process.stdin.on('end', () => {
  try {
    const input = JSON.parse(data);
    const filePath = input.tool_input?.file_path;

    if (filePath && /\.go$/.test(filePath)) {
      const resolvedPath = path.resolve(filePath);
      if (!fs.existsSync(resolvedPath)) { console.log(data); return; }
      // Find nearest go.mod by walking up (max 20 levels to prevent infinite loop)
      let dir = path.dirname(resolvedPath);
      const root = path.parse(dir).root;
      let depth = 0;

      while (dir !== root && depth < 20) {
        if (fs.existsSync(path.join(dir, 'go.mod'))) {
          break;
        }
        dir = path.dirname(dir);
        depth++;
      }

      if (fs.existsSync(path.join(dir, 'go.mod'))) {
        try {
          execFileSync('go', ['vet', './...'], {
            cwd: dir,
            encoding: 'utf8',
            stdio: ['pipe', 'pipe', 'pipe'],
            timeout: 30000
          });
        } catch (err) {
          // go vet exits non-zero when there are errors — filter to edited file
          const output = (err.stdout || '') + (err.stderr || '');
          const relevantLines = output
            .split('\n')
            .filter(line => line.includes(filePath) || line.includes(path.basename(filePath)))
            .slice(0, 10);

          if (relevantLines.length > 0) {
            console.error('[Hook] Go vet errors in ' + path.basename(filePath) + ':');
            relevantLines.forEach(line => console.error(line));
          }
        }
      }
    }
  } catch {
    // Invalid input — pass through
  }

  console.log(data);
});