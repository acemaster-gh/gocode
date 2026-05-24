#!/usr/bin/env node
// =============================================================================
// lib/prompt.js — gocode v1.0.3 UI layer
// @clack/prompts — gocode identity, not a Vite clone
// =============================================================================

import * as p from '@clack/prompts';
import { writeFileSync } from 'fs';

const [,, command, ...args] = process.argv;

// ── ANSI ──────────────────────────────────────────────────────────────────────
const _ = {
    cyan:    s => `\x1b[36m${s}\x1b[0m`,
    green:   s => `\x1b[32m${s}\x1b[0m`,
    yellow:  s => `\x1b[33m${s}\x1b[0m`,
    red:     s => `\x1b[31m${s}\x1b[0m`,
    magenta: s => `\x1b[35m${s}\x1b[0m`,
    dim:     s => `\x1b[2m${s}\x1b[0m`,
    bold:    s => `\x1b[1m${s}\x1b[0m`,
    bg_cyan: s => `\x1b[46m\x1b[30m${s}\x1b[0m`,
};

function finish(data) {
    const f = process.env.GOCODE_RESULT_FILE;
    if (f) writeFileSync(f, JSON.stringify(data));
}

function orExit(val) {
    if (p.isCancel(val)) {
        p.cancel(_.dim('cancelled.'));
        process.exit(1);
    }
    return val;
}

// ── Banner — full ASCII art on wide terminals, compact on narrow ───────────────
function banner(subtitle) {
    const cols = process.stdout.columns || 80;
    console.log("");
    if (cols >= 68) {
        console.log(_.cyan(_.bold("  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗")));
        console.log(_.cyan(_.bold(" ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝")));
        console.log(_.cyan(_.bold(" ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗  ")));
        console.log(_.cyan(_.bold(" ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝  ")));
        console.log(_.cyan(_.bold(" ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗")));
        console.log(_.cyan(_.bold("  ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝")));
        console.log(_.dim("  Developer Workflow System · v1.0.3"));
    } else if (cols >= 40) {
        console.log(_.cyan(_.bold("  ⚡ gocode")) + _.dim(" v1.0.3"));
        console.log(_.dim("  Developer Workflow System"));
    } else {
        console.log(_.cyan(_.bold("⚡ gocode")) + _.dim(" v1.0.3"));
    }
    if (subtitle) console.log(_.dim(`
  ${subtitle}`));
    console.log("");
}

// ── Step label  [2/5] ─────────────────────────────────────────────────────────
function step(n, total, label) {
    return `${_.dim(`[${n}/${total}]`)} ${label}`;
}

// ── Folder tree preview for p.note() ─────────────────────────────────────────
function folderPreview(name, type) {
    const base = `${name}/`;
    const trees = {
        vanilla:  `${base}\n  ├── index.html\n  ├── style.css\n  ├── script.js\n  └── README.md`,
        tailwind: `${base}\n  ├── index.html  ${_.dim('← tailwind CDN')}\n  ├── style.css\n  ├── script.js\n  └── README.md`,
        react:    `${base}\n  ├── src/\n  │   ├── App.jsx\n  │   └── main.jsx\n  ├── index.html\n  ├── vite.config.js\n  └── README.md`,
    };
    return trees[type] || trees.vanilla;
}

// =============================================================================
// CREATE
// =============================================================================
async function cmdCreate() {
    banner('new project');
    p.intro(_.bold(_.bg_cyan(' ⚡ gocode — new project ')));

    // Step 1 — name
    const name = orExit(await p.text({
        message: step(1, 5, 'Project name'),
        placeholder: 'tip-calculator',
        validate: v => !v.trim() ? 'Name cannot be empty.' : undefined
    }));

    // Step 2 — description
    const description = orExit(await p.text({
        message: step(2, 5, 'What does it do?'),
        placeholder: 'e.g. A tab switcher with keyboard navigation',
        initialValue: ''
    }));

    // Step 3 — type
    const type = orExit(await p.select({
        message: step(3, 5, 'Stack'),
        options: [
            { value: 'vanilla',  label: `${_.cyan('⬡')}  Vanilla JS`,   hint: 'HTML · CSS · JS  no build step' },
            { value: 'tailwind', label: `${_.cyan('◈')}  Tailwind CSS`,  hint: 'CDN   utility classes  no build step' },
            { value: 'react',    label: `${_.cyan('⚛')}  React + Vite`,  hint: 'JSX   hot reload   npm scaffold' },
        ]
    }));

    // Step 4 — GitHub
    const create_github = orExit(await p.confirm({
        message: step(4, 5, 'Create GitHub repo?'),
        initialValue: true
    }));

    let visibility = 'public';
    if (create_github) {
        visibility = orExit(await p.select({
            message: _.dim('  └─') + ' Visibility',
            options: [
                { value: 'public',  label: '🌐 Public',  hint: 'visible to everyone' },
                { value: 'private', label: '🔒 Private', hint: 'only you' },
            ]
        }));
    }

    // Step 5 — Netlify
    const deploy_netlify = create_github
        ? orExit(await p.confirm({
            message: step(5, 5, 'Deploy to Netlify?'),
            initialValue: true
          }))
        : false;

    const clean_name = name.trim().toLowerCase().replace(/\s+/g, '-');

    // Preview before build
    p.note(
        [
            _.dim('project') + `  ${_.bold(_.cyan(clean_name))}`,
            _.dim('stack  ') + `  ${type}`,
            _.dim('github ') + `  ${create_github ? visibility : _.dim('skipped')}`,
            _.dim('netlify') + `  ${deploy_netlify ? _.green('yes') : _.dim('skipped')}`,
            '',
            _.dim(folderPreview(clean_name, type)),
        ].join('\n'),
        _.cyan('preview')
    );

    finish({
        name: clean_name,
        description: (description || '').trim(),
        type,
        create_github,
        visibility,
        deploy_netlify
    });
}

// =============================================================================
// OUTRO-CREATE — called by bash after build completes, shows final summary
// args: name type path repo_url netlify_url
// =============================================================================
function cmdOutroCreate() {
    const [name, type, path, repo, netlify] = args;
    const lines = [
        `${_.dim('name   ')}  ${_.bold(name)}`,
        `${_.dim('type   ')}  ${type}`,
        `${_.dim('path   ')}  ${_.dim(path)}`,
    ];
    if (repo && repo !== 'local-only')       lines.push(`${_.dim('github ')}  ${_.cyan(repo)}`);
    if (netlify && netlify !== 'not-deployed') lines.push(`${_.dim('live   ')}  ${_.green(netlify)}`);
    lines.push('');
    lines.push(_.dim('next steps'));
    lines.push(_.dim(`  cd ${path}`));
    if (type === 'react') lines.push(_.dim('  npm run dev'));
    else                  lines.push(_.dim('  open with VS Code Live Server'));
    if (netlify && netlify !== 'not-deployed') lines.push(_.dim(`  visit ${netlify}`));

    p.note(lines.join('\n'), _.green('✓ ready'));
    p.outro(_.bold(_.green('happy coding! ⚡')));
}

// =============================================================================
// RESUME
// =============================================================================
async function cmdResume() {
    const projects = JSON.parse(args[0] || '[]');
    if (projects.length === 0) {
        banner('resume');
        p.log.info('No active projects. Create one: gocode --new');
        process.exit(0);
    }
    banner('resume project');
    p.intro(_.bold(_.bg_cyan(' ⚡ gocode — resume ')));

    const chosen = orExit(await p.select({
        message: 'Pick up where you left off',
        options: projects.map(n => ({ value: n, label: n }))
    }));

    p.outro(_.green(`opening ${_.bold(chosen)}...`));
    finish({ name: chosen });
}

// =============================================================================
// PUSH
// =============================================================================
async function cmdPush() {
    const projectName = args[0] || '';
    banner(`push · ${projectName}`);
    p.intro(_.bold(_.bg_cyan(` ⚡ gocode — push `)));

    const type = orExit(await p.select({
        message: 'Commit type',
        options: [
            { value: 'feat',     label: _.cyan('feat'),     hint: 'new feature' },
            { value: 'fix',      label: _.green('fix'),      hint: 'bug fix' },
            { value: 'style',    label: _.magenta('style'),  hint: 'formatting, UI' },
            { value: 'refactor', label: _.yellow('refactor'),hint: 'restructure' },
            { value: 'docs',     label: 'docs',              hint: 'documentation' },
            { value: 'perf',     label: 'perf',              hint: 'performance' },
            { value: 'test',     label: 'test',              hint: 'tests' },
            { value: 'chore',    label: _.dim('chore'),      hint: 'maintenance' },
        ]
    }));

    const message = orExit(await p.text({
        message: 'What changed?',
        placeholder: 'describe the change clearly',
        validate: v => !v.trim() ? 'Cannot be empty.' : undefined
    }));

    p.outro(_.dim(`${type}: ${message.trim()}`));
    finish({ type, message: message.trim() });
}

// =============================================================================
// DELETE
// =============================================================================
async function cmdDelete() {
    const projects = JSON.parse(args[0] || '[]');
    banner('delete project');
    p.intro(_.bold(_.red(' ⚠  gocode — delete ')));

    const name = orExit(await p.select({
        message: 'Which project to delete?',
        options: projects.map(n => ({ value: n, label: n }))
    }));

    p.log.warn(`This will remove the folder, registry entry, and Obsidian note for ${_.bold(name)}.`);

    const sure = orExit(await p.confirm({
        message: `Permanently delete ${_.bold(_.red(name))}?`,
        initialValue: false
    }));

    if (!sure) { p.cancel(_.dim('nothing deleted.')); process.exit(1); }

    const delete_github = orExit(await p.confirm({
        message: `Also delete the GitHub repo?`,
        initialValue: false
    }));

    p.outro(_.yellow('deleting...'));
    finish({ name, delete_github });
}

// =============================================================================
// EDIT
// =============================================================================
async function cmdEdit() {
    const projects = JSON.parse(args[0] || '[]');
    banner('edit project');
    p.intro(_.bold(_.bg_cyan(' ⚡ gocode — edit ')));

    const name = orExit(await p.select({
        message: 'Which project?',
        options: projects.map(n => ({ value: n, label: n }))
    }));

    const field = orExit(await p.select({
        message: 'What to edit?',
        options: [
            { value: 'name',        label: 'Name',        hint: 'renames folder · GitHub · Obsidian' },
            { value: 'description', label: 'Description',  hint: 'updates GitHub · README · Obsidian' },
            { value: 'visibility',  label: 'Visibility',   hint: 'toggle public / private on GitHub'  },
        ]
    }));

    let value = '';
    if (field === 'name') {
        const raw = orExit(await p.text({
            message: 'New name',
            placeholder: 'e.g. my-new-name',
            validate: v => !v.trim() ? 'Cannot be empty.' : undefined
        }));
        value = raw.trim().toLowerCase().replace(/\s+/g, '-');
    } else if (field === 'description') {
        const raw = orExit(await p.text({
            message: 'New description',
            placeholder: 'e.g. A cool project',
            initialValue: ''
        }));
        value = (raw || '').trim();
    } else if (field === 'visibility') {
        value = orExit(await p.select({
            message: 'Visibility',
            options: [
                { value: 'public',  label: '🌐 Public'  },
                { value: 'private', label: '🔒 Private' },
            ]
        }));
    }

    p.outro(_.green('updating everywhere...'));
    finish({ name, field, value });
}

// =============================================================================
// CONFIG BROWSER — multiselect for browser tabs
// =============================================================================
async function cmdConfigTabs() {
    banner('configure browser tabs');
    p.intro(_.bold(_.bg_cyan(' ⚡ gocode — browser tabs ')));

    const chosen = orExit(await p.multiselect({
        message: 'Which tabs to open on project start?',
        options: [
            { value: 'https://claude.ai',             label: 'Claude AI',       hint: 'claude.ai' },
            { value: 'https://chatgpt.com',            label: 'ChatGPT',         hint: 'chatgpt.com' },
            { value: 'https://developer.mozilla.org',  label: 'MDN Docs',        hint: 'mdn web docs' },
            { value: 'http://localhost:5500',           label: 'Live Server',     hint: 'localhost:5500' },
            { value: 'https://github.com',             label: 'GitHub',          hint: 'github.com' },
            { value: 'https://app.netlify.com',        label: 'Netlify',         hint: 'netlify dashboard' },
            { value: 'https://fonts.google.com',       label: 'Google Fonts',    hint: 'fonts.google.com' },
        ],
        initialValues: [
            'https://claude.ai',
            'https://chatgpt.com',
            'https://developer.mozilla.org',
            'http://localhost:5500',
        ]
    }));

    p.outro(_.green('tabs saved.'));
    finish({ tabs: chosen });
}

// =============================================================================
// LOG HELPERS
// =============================================================================
function cmdLog(level) {
    const msg = args.join(' ');
    switch (level) {
        case 'ok':    p.log.success(msg); break;
        case 'error': p.log.error(msg);   break;
        case 'warn':  p.log.warn(msg);    break;
        case 'info':  p.log.info(msg);    break;
        case 'step':  p.log.step(msg);    break;
        case 'intro': p.intro(msg);       break;
        case 'outro': p.outro(msg);       break;
    }
}

// =============================================================================
// ROUTER
// =============================================================================
switch (command) {
    case 'create':       await cmdCreate();      break;
    case 'resume':       await cmdResume();      break;
    case 'push':         await cmdPush();        break;
    case 'delete':       await cmdDelete();      break;
    case 'edit':         await cmdEdit();        break;
    case 'outro-create': cmdOutroCreate();       break;
    case 'config-tabs':  await cmdConfigTabs();  break;
    case 'ok':
    case 'error':
    case 'warn':
    case 'info':
    case 'step':
    case 'intro':
    case 'outro':        cmdLog(command);        break;
    default:
        process.stderr.write(`unknown: ${command}\n`);
        process.exit(1);
}
