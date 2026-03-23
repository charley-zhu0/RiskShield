#!/usr/bin/env node
/**
 * Stop Hook (Session End) - Persist learnings when session ends
 *
 * Cross-platform (Windows, macOS, Linux)
 *
 * Runs when Claude session ends. Extracts a meaningful summary from
 * the session transcript (via stdin JSON transcript_path) and saves it
 * to a session file for cross-session continuity.
 */

const path = require('path');
const fs = require('fs');
const {
  getSessionsDir,
  getDateString,
  getTimeString,
  getSessionIdShort,
  ensureDir,
  readFile,
  writeFile,
  replaceInFile,
  findFiles,
  log
} = require('../lib/utils');

/**
 * Extract a meaningful summary from the session transcript.
 * Reads the JSONL transcript and pulls out key information:
 * - User messages (tasks requested) → mapped to Completed items
 * - Tools used
 * - Files modified
 *
 * Claude Code transcript format (JSONL, two supported layouts):
 *
 * Layout A (Claude Code native JSONL):
 *   { "type": "user", "message": { "role": "user", "content": string|array }, ... }
 *   { "type": "assistant", "message": { "role": "assistant", "content": array }, ... }
 *   Tool use blocks live inside message.content: [{ "type": "tool_use", "name": "Edit", "input": {...} }]
 *
 * Layout B (simple flat format):
 *   { "role": "user"|"assistant", "content": string|array }
 *   Tool use blocks are also in content array.
 */
function extractSessionSummary(transcriptPath) {
  const content = readFile(transcriptPath);
  if (!content) return null;

  const lines = content.split('\n').filter(Boolean);
  const userMessages = [];
  const toolsUsed = new Set();
  const filesModified = new Set();
  let parseErrors = 0;

  for (const line of lines) {
    try {
      const entry = JSON.parse(line);

      // ---- Resolve message object (handles both Layout A and B) ----
      // Layout A: entry.type === 'user'/'assistant', actual message in entry.message
      // Layout B: entry.role === 'user'/'assistant', content directly on entry
      const msg = (entry.type === 'user' || entry.type === 'assistant')
        ? entry.message   // Layout A
        : entry;          // Layout B

      if (!msg) continue;

      // ---- Collect user messages (real human prompts only) ----
      if (msg.role === 'user') {
        let text = '';
        if (typeof msg.content === 'string') {
          text = msg.content;
        } else if (Array.isArray(msg.content)) {
          // Extract only text blocks; skip tool_result blocks (system injected)
          text = msg.content
            .filter(c => c && c.type === 'text')
            .map(c => c.text || '')
            .join(' ');
        }
        const cleaned = text.replace(/<[^>]+>/g, '').trim();
        if (cleaned && cleaned.length > 5) {
          userMessages.push(cleaned.slice(0, 300));
        }
      }

      // ---- Collect tools used and files modified ----
      if (msg.role === 'assistant' && Array.isArray(msg.content)) {
        for (const block of msg.content) {
          if (block && block.type === 'tool_use') {
            const toolName = block.name || '';
            if (toolName) toolsUsed.add(toolName);

            const filePath = block.input?.file_path || '';
            if (filePath && (toolName === 'Edit' || toolName === 'Write')) {
              filesModified.add(filePath);
            }
          }
        }
      }

      // ---- Also handle flat tool_use entries ----
      if (entry.type === 'tool_use') {
        const toolName = entry.name || entry.tool_name || '';
        if (toolName) toolsUsed.add(toolName);
        const filePath = entry.input?.file_path || entry.tool_input?.file_path || '';
        if (filePath && (toolName === 'Edit' || toolName === 'Write')) {
          filesModified.add(filePath);
        }
      }
    } catch {
      parseErrors++;
    }
  }

  if (parseErrors > 0) {
    log(`[SessionEnd] Skipped ${parseErrors}/${lines.length} unparseable transcript lines`);
  }

  if (userMessages.length === 0 && toolsUsed.size === 0) return null;

  return {
    userMessages: userMessages.slice(-10), // Last 10 user messages = tasks worked on
    toolsUsed: Array.from(toolsUsed).slice(0, 20),
    filesModified: Array.from(filesModified).slice(0, 30),
    totalMessages: userMessages.length
  };
}

// Read hook input from stdin (Claude Code provides transcript_path via stdin JSON)
const MAX_STDIN = 1024 * 1024;
let stdinData = '';
process.stdin.setEncoding('utf8');

process.stdin.on('data', chunk => {
  if (stdinData.length < MAX_STDIN) {
    stdinData += chunk;
  }
});

process.stdin.on('end', () => {
  runMain();
});

function runMain() {
  main().catch(err => {
    console.error('[SessionEnd] Error:', err.message);
    process.exit(0);
  });
}

async function main() {
  // Parse stdin JSON to get transcript_path
  let transcriptPath = null;
  try {
    const input = JSON.parse(stdinData);
    transcriptPath = input.transcript_path;
  } catch {
    // Fallback: try env var for backwards compatibility
    transcriptPath = process.env.CLAUDE_TRANSCRIPT_PATH;
  }

  const sessionsDir = getSessionsDir();
  const today = getDateString();
  const shortId = getSessionIdShort();

  ensureDir(sessionsDir);

  const currentTime = getTimeString();

  // Try to extract summary from transcript
  let summary = null;

  if (transcriptPath) {
    if (fs.existsSync(transcriptPath)) {
      summary = extractSessionSummary(transcriptPath);
      if (!summary) {
        log('[SessionEnd] Could not extract summary from transcript (empty or unrecognized format)');
      }
    } else {
      log(`[SessionEnd] Transcript not found: ${transcriptPath}`);
    }
  } else {
    log('[SessionEnd] No transcript_path provided in stdin');
  }

  // Find existing session file for today (any short-id for today's date)
  const existingSessions = findFiles(sessionsDir, `${today}-*-session.tmp`);
  // Also check old format: YYYY-MM-DD-session.tmp
  const oldFormatFile = path.join(sessionsDir, `${today}-session.tmp`);
  const oldFormatExists = fs.existsSync(oldFormatFile);

  // Prefer existing file for today; otherwise use shortId-based filename
  let sessionFile;
  if (existingSessions.length > 0) {
    sessionFile = existingSessions[0].path;
  } else if (oldFormatExists) {
    sessionFile = oldFormatFile;
  } else {
    sessionFile = path.join(sessionsDir, `${today}-${shortId}-session.tmp`);
  }

  if (fs.existsSync(sessionFile)) {
    // Update timestamp
    replaceInFile(
      sessionFile,
      /\*\*Last Updated:\*\*.*/,
      `**Last Updated:** ${currentTime}`
    );

    // Replace blank template with real summary if we have one
    if (summary) {
      const existing = readFile(sessionFile);
      if (existing && existing.includes('[Session context goes here]')) {
        // Replace the entire Current State block with summary
        const updatedContent = existing.replace(
          /## Current State[\s\S]*$/,
          buildSummarySection(summary) + '\n'
        );
        writeFile(sessionFile, updatedContent);
      } else if (existing) {
        // Session already has content — append a new summary block
        const appendContent = `\n---\n\n## Session Update (${today} ${currentTime})\n\n${buildSummarySection(summary)}\n`;
        fs.appendFileSync(sessionFile, appendContent, 'utf8');
      }
    }

    log(`[SessionEnd] Updated session file: ${sessionFile}`);
  } else {
    // Create new session file
    const summarySection = summary
      ? buildSummarySection(summary)
      : `## Current State\n\n[Session context goes here]\n\n### Completed\n- [ ]\n\n### In Progress\n- [ ]\n\n### Notes for Next Session\n-\n\n### Context to Load\n\`\`\`\n[relevant files]\n\`\`\``;

    const template = `# Session: ${today}
**Date:** ${today}
**Started:** ${currentTime}
**Last Updated:** ${currentTime}

---

${summarySection}
`;

    writeFile(sessionFile, template);
    log(`[SessionEnd] Created session file: ${sessionFile}`);
  }

  process.exit(0);
}

/**
 * Build the summary section of the session file.
 * Uses Completed/In Progress format compatible with session-manager.js parseSessionMetadata.
 */
function buildSummarySection(summary) {
  let section = '## Current State\n\n';

  // Completed tasks — treat last N user messages as tasks that were worked on
  section += '### Completed\n';
  if (summary.filesModified.length > 0) {
    // If files were modified, list them as completed work items
    for (const f of summary.filesModified) {
      section += `- [x] Modified: ${path.basename(f)}\n`;
    }
  } else if (summary.userMessages.length > 0) {
    // Fall back to user messages as task descriptions
    for (const msg of summary.userMessages.slice(0, 5)) {
      section += `- [x] ${msg.replace(/`/g, '\\`').split('\n')[0].slice(0, 120)}\n`;
    }
  } else {
    section += '- [ ]\n';
  }
  section += '\n';

  // In Progress — list any tasks beyond the completed ones
  section += '### In Progress\n';
  const inProgressMessages = summary.userMessages.slice(5);
  if (inProgressMessages.length > 0) {
    for (const msg of inProgressMessages) {
      section += `- [ ] ${msg.replace(/`/g, '\\`').split('\n')[0].slice(0, 120)}\n`;
    }
  } else {
    section += '- [ ]\n';
  }
  section += '\n';

  // Notes — tools used as context
  section += '### Notes for Next Session\n';
  if (summary.toolsUsed.length > 0) {
    section += `- Tools used: ${summary.toolsUsed.join(', ')}\n`;
  }
  section += `- Total user messages: ${summary.totalMessages}\n`;
  section += '\n';

  // Context to load — files modified
  section += '### Context to Load\n```\n';
  if (summary.filesModified.length > 0) {
    section += summary.filesModified.join('\n') + '\n';
  } else {
    section += '[relevant files]\n';
  }
  section += '```';

  return section;
}
