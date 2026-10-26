| Release | Date       | Comments                                                             |
|---------|------------|----------------------------------------------------------------------|
| 2.9.1   | 2026.10.26 | **Breaking:** `policy backup`/`policy restore` now default to a single named policy (`policy backup POLICY_NAME FILE`); the v2.9.0 all-policies-in-one-file behavior moved behind `--all`. `policy backup --all` also now skips Vault's built-in/virtual policies (e.g. `root`, which `ListPolicies` reports but has no stored rules to read back) instead of aborting the whole backup. `show` moved from a `policy list` alias to a `policy read`/`get` alias. |
| 2.9.0   | 2026.10.26 | Added `policy backup`/`policy restore`, dumping and restoring every ACL policy in a single JSON or HCL file; fixed `kv write --batch`/`--json` failing on a path whose latest version was soft-deleted or destroyed; a typo'd subcommand (e.g. `kv bogus`) now correctly exits non-zero instead of silently exiting 0 |
| 2.8.0   | 2026.09.27 | `kv write --json` writes a whole secret from a flat JSON file in a single call; every command's `--help` now includes runnable examples |
| 2.7.0   | 2026.09.27 | `kv write --batch` writes multiple KEY/VALUE fields from a file in a single call |
| 2.6.0   | 2026.09.22 | `kv write` can now read a secret from a file and inject it in Vault |
| 2.5.0   | 2026.09.22 | `kv read` can now write secret to a file |
| 2.4.3   | 2026.08.07 | Errors now go to stderr, not stdout                                  | 
| 2.4.2   | 2026.08.06 | Minor TUI fix in `kv write`                                          | 
| 2.4.1   | 2026.07.06 | Refactored code to make it more _GO idiomatic_                       |
| 2.4.0   | 2026.07.04 | Completed the `token` subcommand; software is now feature-complete   |
| 2.3.1   | 2026.06.28 | Backed up engines now dump their file encoded (option for cleartext) |
| 2.3.0   | 2026.06.27 | Added kv engine backup and restore functions                         |
| 2.2.0   | 2026.06.26 | Renamed the `secrets` command as `kv`, created the `sys` command     |
| 2.1.0   | 2026.06.21 | Now aligned with SemVer<br>Policies and tokens subcommands           |
| 2.00.00 | 2026.06.18 | Initial version with `secrets` and `admin` subcommands               |
