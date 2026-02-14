#!/usr/bin/env node
/**
 * PostToolUse Hook: Auto-format Go files with goimports after edits
 *
 * Cross-platform (Windows, macOS, Linux)
 *
 * Runs after Edit tool use. If the edited file is a Go file,
 * formats it with goimports. Fails silently if goimports isn't installed.
 */

const { execFileSync } = require('child_process');

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
      try {
        execFileSync('goimports', ['-w', filePath], {
          stdio: ['pipe', 'pipe', 'pipe'],
          timeout: 15000
        });
      } catch {
        // goimports not installed, file missing, or failed — non-blocking
      }
    }
  } catch {
    // Invalid input — pass through
  }

  console.log(data);
});