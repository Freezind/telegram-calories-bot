用cursor或claude code及同类产品，95%以上通过vibe coding的方式实现一个telegram bot、

产品功能：

- 用户上传图片后计算卡路里（LUI）
- 可以查询历史的数据列表（MINI APP）

需要包含：

- LUI 部分
    - /slash 命令
    - 通过 LUI 及 Inline button实现一套bot的交互
- Mini app 部分
    - GUI 实现数据列表的 CRUD
    - 通过LLM 及相关任一工具进行的 bot 测试

交付物

- 可使用的bot name
- 可以看的代码（需要用有效的方式控制claude code的代码质量）
- vibe coding 的完整提示词
- llm 测试工具的提示词
- 输出的测试报告

技术要求
- 使用golang, 和 https://github.com/tucnak/telebot 框架