# Notes about how to config zed

### tasks.json

```
[
  {
    "label": "Git Commit All (VS Code style)",
    "command": "git add -A; if ($?) { $msg = Read-Host 'Enter commit message'; git commit -m $msg }",
    "use_new_terminal": false,
    "allow_concurrent_runs": false,
  },
]

```

### .zed/debug.json

```
[
  {
    "label": "Launch API Root",
    "adapter": "Delve",
    "request": "launch",
    "mode": "debug",
    "program": "./cmd/api",
    "cwd": "$ZED_WORKTREE_ROOT",
  },
]

```

### keymap.json

```
[
  {
    "context": "Workspace",
    "bindings": {
      "alt-i": [
        "task::Spawn",
        {
          "task_name": "Git Commit All (VS Code style)",
        },
      ],
    },
  },
]

```
