# TODO LIST
___

| Task                                           | Slated for   | Comments |
|------------------------------------------------|--------------|----------|
| [x] secret.Read()                              | 2.0.0-debug1 |          |
| [x] secret.Write()                             | 2.0.0-debug2 |          |
| [x] secret.Delete()                            | 2.0.0-debug3 |          |
| [x] secret.Destroy()                           | 2.0.0-debug3 |          |
| [x] admin.SetRootKeys()                        | 2.0.0-debug4 |          |
| [x] admin.Seal()                               | 2.0.0-debug4 |          |
| [x] admin.Unseal()                             | 2.0.0-debug4 |          |
| [x] policy list                                | 2.1.0        |          |
| [ ] policy create                              |              |          |
| [x] policy delete                              | 2.1.0        |          |
| [ ] policy edit                                |              |          |
| [x] list mounts                                | 2.1.0        |          |
| [ ] token create                               |              |          |
| [ ] token delete                               |              |          |
| [ ] token edit                                 |              |          |
| [ ] better progress handling in admin.Unseal() |              |          |


___
# TASKS

- [ ] Allow `admin setrootkeys` to work in offline mode : 
  - bypass admin.GetSealStatus()
  - the `MinimumRequired` member of `VaultRootKeysStruct` should be set to zero in the meantime