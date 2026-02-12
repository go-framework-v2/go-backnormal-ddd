# LoginService 的 db 与仓储注入说明

本文档说明两件事：**（1）Service 里的 `db` 是什么、和事务 `tx` 的关系、以及未来加缓存等如何扩展；（2）为什么还要传入 `appRepo`、`userRepo`，它们有什么用。**

---

## 一、`db` 是什么？是 tx 吗？

### 1.1 结论：不是

- **`db`**：注入到 LoginService 里的那个字段，表示**应用层可用的「数据库入口」**。  
  在当前项目里就是 `res.MysqlDB`，即**主库的连接池入口**（一个 `*gorm.DB`）。  
  可以理解为：**“这个应用能用来访问数据库的那把钥匙”**，和“当前有没有开事务”无关。

- **`tx`**：是**某一次业务里开启的一笔事务**。  
  只有在代码里执行 `db.Transaction(func(tx *gorm.DB) error { ... })` 时，回调里拿到的 `tx` 才是“当前这笔事务”。  
  所以：**tx 是“某次请求里、某一段操作”的数据库连接，db 是“应用全局的、用来开事务/查库”的入口。**

### 1.2 关系一句话

- **`db`**：应用层持有的、用来“开事务”或“直接查库”的数据库连接（你现在就是主库）。
- **`tx`**：由 `db.Transaction(...)` 在回调里给你的**当前这笔事务**；只在这次回调内有效，用来保证这一段里的多次 DB 操作原子。

所以：**“db 表示应用可操作的数据库连接”是对的；tx 是“当前这一次事务”，由 db 在 Transaction 里创建，不是 Service 的字段。**

---

## 二、为什么要给 LoginService 加一个 `db` 字段？

目的只有一个：**让应用层能控制事务边界。**

- 如果 Service 里**没有** `db`，只有 `appRepo`、`userRepo`（且它们是用“某次创建时的 db”构造的），那么你没法在应用层说：**“下面这几步（查 App、查/插 User、更新 User）必须在同一笔事务里”**。事务只能写在 repo 内部，这样「查 App + FindByUk + Update」就跨了多笔连接/多次提交，**不能原子提交或回滚**。

- 给 Service 注入 **`db`** 之后，应用层可以写：`s.db.Transaction(func(tx *gorm.DB) error { ... })`，在回调里用 **`tx`** 再创建一遍 repo（`NewBizAppRepository(tx)`、`NewBizUserRepository(tx)`），这样**查 App、FindByUk（查或插）、UpdateByFieldmap** 都用同一个 `tx`，要么一起提交，要么一起回滚。

所以：**`db` 表示“应用对象可以操作的数据库连接”，用来在应用层开事务；tx 是某次 `Transaction` 里的事务连接，由 db 创建出来。**

---

## 三、未来加缓存、只读库等怎么扩展？

思路是：**Service 里继续放“应用层可用的各种资源入口”，db 只是其中“可写的主库”这一种。**

| 扩展       | 做法 |
|------------|------|
| **加缓存（如 Redis）** | 在 LoginService 里再加字段，例如 `cache *redis.Client` 或 `cache CacheInterface`。读用户时先查 cache，未命中再查 db；**写用户**仍然用 `db.Transaction` 包住「查 App + 查/插/改 User」。事务只包 DB 写路径；缓存在事务提交后再更新或异步失效。 |
| **加只读从库（读从、写主）** | 可以有两个字段：`dbWrite *gorm.DB`（主库）、`dbRead *gorm.DB`（从库）。像 GuestLogin 这种要查+可能要插/改的流程，**整段仍然用 `dbWrite.Transaction`**，保证读写都在主库、同一笔事务。纯查询接口可以单独用 `dbRead` 或基于它建的 repo，减轻主库压力。 |
| **加 MQ、外部 RPC 等** | 同样：Service 里多一个 `mq` 或 `eventPublisher` 等字段。事务里只做 DB 操作；事务成功提交之后，再发 MQ/事件，避免“消息发了但 DB 回滚”的不一致。 |

**总结**：**`db` 表示“应用对象可以操作的（主）数据库连接”，用来开事务（tx）；未来加缓存、只读库、MQ 等，都是在 Service 上再加别的“资源入口”，事务仍然只围绕 `db`（或 dbWrite）来开。**

---

## 四、`appRepo`、`userRepo` 有什么用？

### 4.1 当前事实

- **GuestLogin** 里：只用了 **`s.db`**（用来 `s.db.Transaction(...)`）；事务回调里用的是 **`repository.NewBizAppRepository(tx)`、`repository.NewBizUserRepository(tx)`**，没有用 `s.appRepo`、`s.userRepo`。
- 所以**仅就当前这一个接口而言**，注入的 `appRepo`、`userRepo` 确实**没有被用到**；有用的是传进来的 `db`。

### 4.2 为什么还要传？

主要是为**后续扩展**和**单测**预留，而不是为当前这一条 GuestLogin 服务：

| 用途 | 说明 |
|------|------|
| **以后「只读、不需要事务」的接口** | 例如「根据 projectId 查 App」「根据 userId 查用户」等，直接用 `s.appRepo.FindByProjectID(...)`、`s.userRepo.FindByXXX(...)` 即可，**不需要**在方法里再写 `db.Transaction`，也不必每次 `NewXxxRepository(db)`。 |
| **单测 / Mock** | 测试 LoginService 时，可以注入 **mock 的** `AppRepository`、`UserRepository`，不依赖真实 DB 和事务，用例更好写、运行更快。 |
| **依赖一眼能看懂** | 构造函数是 `NewLoginService(db, appRepo, userRepo)`，能直接看出：这个 Service 会用到「数据库 + App 仓储 + User 仓储」，职责边界清晰。 |

所以：**`db` 负责「开事务」；`appRepo` / `userRepo` 负责「在不需要事务的场景里直接用、以及方便测试和扩展」。** 当前只实现了 GuestLogin，所以看起来像是“多传了两个没用的”，但从整体设计上保留是合理的。

### 4.3 如果坚持「不传没用到的」

也可以改成**只传 `db`**，在需要事务的地方再临时 New：

```go
// 只注入 db
loginService = service.NewLoginService(db)

// GuestLogin 里仍然 s.db.Transaction，里面 NewBizAppRepository(tx)、NewBizUserRepository(tx)
// 以后若有「只查 App」的接口，再在 Service 里 NewBizAppRepository(s.db) 用一次
```

这样做的代价是：以后每个「不需要事务」的接口里可能都要写一遍 `repository.NewBizXxxRepository(s.db)`，或者再改回给 Service 注入 repo。所以**更常见的做法是一开始就注入 `db + appRepo + userRepo`**，哪怕当前只用到 `db`。

---

## 五、总结

| 参数 / 字段 | 含义 | 在当前 GuestLogin 里的作用 |
|-------------|------|----------------------------|
| **db** | 应用可用的数据库连接（主库入口） | 用来 `s.db.Transaction(...)` 开事务 |
| **appRepo** | 基于某次传入的 db 创建的 App 仓储 | 当前未用；留给以后只读接口和单测 |
| **userRepo** | 基于某次传入的 db 创建的 User 仓储 | 当前未用；留给以后只读接口和单测 |
| **tx**（非字段） | 某次 `Transaction` 回调里的事务连接 | 在 GuestLogin 的回调里用 tx 创建 `appRepoTx`、`userRepoTx`，保证整段原子 |

**一句话**：**db 表示“应用对象可以操作的数据库连接”，用来开事务；appRepo/userRepo 为“非事务场景 + 单测”预留；tx 是某次事务的连接，由 db 在 Transaction 里创建。**
