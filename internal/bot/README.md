# bot

Используеться пакет `github.com/PaulSonOfLars/gotgbot/v2`, сами разработчики утверждают следющее:

1. All telegram API types and methods are generated from the bot api docs, which makes this library:
    - Guaranteed to match the docs
    - Easy to update
    - Self-documenting (Re-uses pre-existing telegram docs)
2. Type safe; no weird interface{} logic, all types match the bot API docs.
3. No third party library bloat; only uses standard library.
4. Updates are each processed in their own go routine, encouraging concurrent processing, and keeping your bot responsive.
5. Code panics are automatically recovered from and logged, avoiding unexpected downtime.

## Commands

 - `/find`
 - `/ping`