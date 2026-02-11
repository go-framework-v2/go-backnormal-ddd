# DDD 落地实践指南 —— 以图书馆系统为例

> 本文结合本项目实际代码，通俗讲解 DDD 如何落地、用了哪些方法论，以及新业务/改动时如何操作。

---

## 一、DDD 是什么？一句话理解

**DDD（领域驱动设计）的核心思想：让代码结构围绕「业务领域」组织，而不是围绕「技术实现」。**

传统做法：按技术分层（Controller → Service → DAO），业务逻辑散落各处。
DDD 做法：先识别业务中的核心概念（用户、借阅、书籍），让它们成为代码的「主角」，技术只是配角。

---

## 二、本项目的分层架构

```
┌─────────────────────────────────────────────────────────────────┐
│  接口层 (api/)                                                    │
│  - 接收 HTTP 请求，参数校验                                        │
│  - 调用应用服务，返回响应                                           │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  应用层 (internal/application/)                                   │
│  - Application Service：编排业务流程，不包含核心业务规则              │
│  - DTO：请求/响应的数据结构                                         │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  领域层 (internal/domain/)           ← 核心！业务规则集中在这里        │
│  - 聚合根、实体、值对象、领域事件                                    │
│  - 仓储接口（只定义，不实现）                                        │
└──────────────────────────┬──────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│  基础设施层 (internal/infrastructure/)                             │
│  - 仓储实现：把领域对象存到 MySQL                                    │
│  - PO（持久化对象）：表结构对应的模型                                 │
└─────────────────────────────────────────────────────────────────┘
```

**依赖方向**：接口层 → 应用层 → 领域层 ← 基础设施层
领域层**不依赖**任何外部技术，只依赖业务概念。

---

## 三、用到的 DDD 核心方法论

### 1.  bounded context（限界上下文）

把大系统拆成多个「业务子域」，每个子域有自己的一套模型和术语。

在本项目中：


| 限界上下文            | 职责                     | 主要概念                      |
| --------------------- | ------------------------ | ----------------------------- |
| **用户上下文** (user) | 用户注册、信息、地址     | User, UserAddress, UserStatus |
| **书籍上下文** (book) | 书籍管理、库存           | Book, BookID                  |
| **借阅上下文** (loan) | 借书、还书、续借、逾期费 | Loan, LoanStatus              |

不同上下文之间通过 **ID** 引用，而不是直接持有对象。
例如：`Loan` 只持有 `user.UserID`、`book.BookID`，不直接持有 `*User`、`*Book`。

---

### 2.  aggregate root（聚合根）

聚合是「一组强相关的对象的集合」，聚合根是这组对象的入口，外界只能通过聚合根修改内部状态。

**本项目的聚合根**：

- **User**：用户聚合根，管理 `UserAddress` 等
- **Book**：书籍聚合根
- **Loan**：借阅聚合根

原则：

- 业务规则写在聚合根里（如 `User.CanBorrow()`、`Loan.Return()`）
- 修改聚合内部状态必须通过聚合根的方法
- 一个事务通常只改一个聚合根

示例：

```go
// 借书：通过 User 聚合根执行业务规则
func (u *User) BorrowBook(bookID book.BookID) error {
    if err := u.CanBorrow(); err != nil {  // 业务规则
        return err
    }
    u.currentLoans++
    u.updatedAt = time.Now()
    // ...
}
```

---

### 3.  entity（实体）vs value object（值对象）


| 类型       | 特点                                       | 本项目中例子                               |
| ---------- | ------------------------------------------ | ------------------------------------------ |
| **实体**   | 有唯一标识，生命周期内会变化               | User, Book, Loan, UserAddress              |
| **值对象** | 无标识，用属性是否相同判断相等，通常不可变 | Address, Money, UserID, BookID, LoanStatus |

**值对象的好处**：

- 把校验和业务规则封装在创建时（如 `NewAddress` 校验省市、电话）
- 不可变，避免被意外修改
- 可以在多个聚合中复用（如 `Address`）

示例：

```go
// 值对象：创建时校验，创建后不可变
addr, err := valueobject.NewAddress(province, city, recipient, phone)
if err != nil {
    return nil, err  // 无效地址直接拒绝
}

// 借阅时保存地址快照，即使用户 later 改了地址，历史借阅记录不受影响
addressJSON, _ := deliveryAddress.ToJSON()
loan.deliveryAddress = addressJSON
```

---

### 4.  repository（仓储）

仓储是「领域对象的持久化抽象」：领域层只定义接口，基础设施层实现。

```go
// 领域层：只定义接口
type UserRepository interface {
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// 基础设施层：实现接口，操作 MySQL
type UserRepositoryImpl struct {
    db *gorm.DB
}
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
    // 查 DB -> 转成领域对象 User
}
```

这样领域层不依赖 GORM、MySQL，测试时可以轻松用 Mock 替换。

---

### 5.  factory（工厂）

复杂对象的创建逻辑封装在工厂方法里，保证创建出来的对象始终合法。

```go
// 工厂方法：校验 + 构造
func NewUser(username, email, password string) (*User, error) {
    if username == "" {
        return nil, errors.New("username is required")
    }
    // ...
    return &User{...}, nil
}

// 从持久化还原：Repository 内部使用
func RestoreUser(id UserID, username string, ...) *User {
    return &User{...}  // 已知来自 DB，不做完整校验
}
```

---

### 6.  domain events（领域事件）

聚合内部发生重要业务变化时，记录「领域事件」，便于后续扩展（如发通知、统计、同步）。

```go
// Loan 归还时发出事件
l.addDomainEvent(LoanReturnedEvent{
    LoanID:     l.id,
    UserID:     l.userId,
    BookID:     l.bookId,
    ReturnedAt: now,
    LateFee:    l.lateFee,
    Timestamp:  now,
})

// 应用层或基础设施层可订阅这些事件，做后续处理
events := loan.ClearDomainEvents()
for _, e := range events {
    // 发 MQ、写日志、更新统计等
}
```

目前项目主要做了「记录」，事件发布可以后续接入 MQ 等。

---

## 四、一次请求的完整链路（以借书为例）

```
HTTP POST /api/ddd/library/loans
    ↓
[api] handleBorrowBook：绑定 JSON → BorrowBookReq
    ↓
[application] LibraryService.BorrowBook：
    1. 用 userRepo.FindByID 查出 User
    2. 用 bookRepo.FindByID 查出 Book
    3. valueobject.NewAddress 构造地址
    4. u.CanBorrow()、b.CanBorrow() 做业务校验
    5. loan.NewLoan 创建借阅
    6. u.BorrowBook()、b.Borrow() 更新聚合
    7. userRepo.Save、bookRepo.Save、loanRepo.Save 持久化
    ↓
[domain] User.BorrowBook、Book.Borrow、Loan 创建
    - 所有业务规则在这里执行
    ↓
[infrastructure] UserRepositoryImpl.Save
    - User → UserPO，写入 ddd_user、ddd_user_address
```

**要点**：
应用服务负责「流程编排」，领域对象负责「业务规则」，仓储负责「持久化」。

---

## 五、新业务 / 改动时的操作指南

### 场景 1：新增一个「业务能力」（如：书籍预约）

**步骤**：

1. **领域层**

   - 判断是否新聚合：预约可视为 `Reservation` 聚合根
   - 新建 `internal/domain/reservation/`
   - 定义 `Reservation`、`ReservationRepository` 接口
   - 把预约规则写在 `Reservation` 里（如最多预约几本、预约有效期）
2. **基础设施层**

   - 建表 `ddd_reservation`
   - 新建 `ReservationPO`
   - 实现 `ReservationRepositoryImpl`
3. **应用层**

   - 新建 `ReserveBookReq/Resp` DTO
   - 在 `LibraryService` 中加 `ReserveBook` 方法，调用 `Reservation` 和 Repository
4. **接口层**

   - 新增 `POST /api/ddd/library/reservations`
   - 在 route 中绑定 handler

---

### 场景 2：修改已有业务规则（如：逾期费从每天 0.5 元改为 1 元）

**只改领域层**：

```go
// internal/domain/loan/loan.go
func (l *Loan) calculateLateFee(overdueDays int) float64 {
    if overdueDays <= 3 {
        return 0
    }
    baseFee := 2.0
    dailyFee := 1.0  // 从 0.5 改为 1.0
    return baseFee + (float64(overdueDays-3) * dailyFee)
}
```

接口层、应用层、基础设施层都**不用动**。

---

### 场景 3：给 User 增加一个字段（如：实名认证状态）

1. **领域层**：在 `User` 中加字段和访问方法
2. **基础设施层**：在 `UserPO` 和表结构中加列，更新 `toDomain` / `toPO`
3. **应用层**：如需对外暴露，在 DTO 中增加字段

---

### 场景 4：换数据库（如 MySQL → PostgreSQL）

只改基础设施层：

- 新建 `internal/infrastructure/persistence/postgres/`
- 实现各 `XxxRepository` 接口
- 在启动/注入处替换为 Postgres 版本

领域层、应用层、接口层都不需要改动。

---

## 六、实战 checklist


| 操作            | 建议                                                        |
| --------------- | ----------------------------------------------------------- |
| 加业务规则      | 先想清楚属于哪个聚合，写在对应聚合根里                      |
| 加新实体/值对象 | 先判断是实体还是值对象，再确定放在哪个上下文                |
| 跨聚合操作      | 在应用服务里编排，分别调用多个仓储和聚合根                  |
| 加新接口        | 接口层加 handler → 应用层加 service 方法 → 需要时改领域层 |
| 换存储/中间件   | 只改基础设施层，通过实现/替换接口完成                       |

---

## 七、总结

1. **分层清晰**：接口 → 应用 → 领域 ← 基础设施，领域是核心。
2. **以领域为中心**：业务规则在聚合根和值对象里，应用层只做编排。
3. **依赖倒置**：领域定义仓储接口，基础设施实现，便于测试和替换。
4. **改动局部化**：改规则动领域层，改流程动应用层，改存储动基础设施层。

按照上述方式拆业务、写领域、做持久化，就能在实战中持续落地 DDD。
