BartmossRoar 设计记忆点

1. 身份与密码学

· 长期身份：SPHINCS+ Small 签名密钥，公钥哈希即 ID，强绑定。
· 密钥协商：每 30min 轮换 ML-KEM 封装密钥；每 10min 重协商 XChaCha20Poly1305 对称密钥。
· 协商保护：协商包由上一次对称密钥加密，混入洪流广撒网，不与普通消息区分。
· 首条消息：对方需在发送第一条信息前触发密钥协商。
· PoT 反馈签名：使用 epoch 一次性密钥签名，不暴露长期身份。

2. 洪流与网络层

· 发送频率：3-7s 随机 tick 发送固定大小 block（无论有无消息，全发 dummy）。
· 匿名集：对方所在 channel 内所有成员。channel 成员越多越匿名。
· Channel 派生：基于 secret + 时间的伪随机一致性计算，约 60min 轮换一次。
· Channel 细分：每个成员只向自己此刻的 channel 发掩盖流量，不盖满全网。
· 传输层：选用 tornago（Go 封装 C 版 Tor），通过 .onion 隐藏服务解决严格 NAT，不需要 Gonc（互斥关系，不混用）。
· Tag：数据包附带 tag，包含可能接收对象。dummy 也有 tag，且带随机长期“伪装好友”常驻 tag 内，加快传播但模糊路径。

3. PoT（Proof of Transfer）与抗女巫

· 发送权：节点需积攒一定数量的 PoT 才能发送消息。
· 转发机制：帮别人转发消息，Head 部分签名替换为自己对 Payload 的签名，Payload 解密后包含 FromID 和 FromID 对 Payload 的签名（只有接收方知道发送者是谁）。
· 反馈与收集：反馈不走洪流，不必须匿名（观察谁在攒 PoT 无效益）。接收方通过带外通道回签“承认你在时间 T 转发了消息 ID X”，节点收集这些签名形成 PoT。
· 窗口限制：PoT 保留 60min，超出不承认。
· 抗女巫逻辑：邻居每小时随机轮换，攻击者很难刚好刷出一个 SPHINCS+ 公钥哈希是现有女巫节点的邻居。

4. 存储（BadgerDB v4）—— 6 张表

done · table0：已知节点的 ID 和 pubkey（SPHINCS+ Small）。
done · table1：channel 内成员的 MLDSA 公钥和已交换的 XChaCha 密钥（30min/10min 定时重发与协商）。
done · table2：10s TTL，保存已解密/已广播消息头的 Blake3 哈希（防重放）。
done · table3（可选）：消息记录（id、from、timestamp、content、可选原数据包压缩 bytes）。
· table4（新增）：记录 PoT 签名（带外反馈回来的签名，绑定 epoch 一次性密钥和消息 ID X）。
done · table5（新增）：60min 内的消息 ID（用于 PoT 验证与防伪造，超出 60min 的消息 ID 不承认）。

5. 开发顺序（极简 Prototype）

done 6. 写 main 和顶层函数，定好接口。
done 7. 定 transport 接口（Send/OnRecv）。
done 8. 先用 TCP 直连跑通全链路。
done 9. 搞 BadgerDB 前四张表（CRUD + TTL）。
10. 写洪流 tick（3-7s，先全 dummy）。
11. 加 tag 和头部解密（区分“给我的”和“不是给我的”）。
12. 接 MLKEM + XChaCha 协商（先手动，后定时）。
13. 加 channel 发现（先固定名单，再 secret+时间派生）。
done 14. 把 tornago 塞进 transport，替换 TCP，别动上层。
15. 两台机器 + 一个 VPS 真跑，调 NAT、掉线、重连。
16. 实现 PoT 基础逻辑：转发签名替换、带外反馈通道、收集签名。
17. 引入 table4 和 table5：存储 PoT 签名、维护 60min 消息 ID 窗口。
18. 加 PoT 发送配额验证：无足够 PoT 则拒绝发送/转发。
19. 最后补消息记录、压缩、优化。

