# test_task

Задание по которому реализовывал сервис:
Design and implement "Word of Wisdom" tcp server.  
• TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Proof Of Work verification, server should send one of the quotes from "word of wisdom" book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

## Выбор PoW алгоритма

В качестве PoW алгоритма используется HashCash-подобный алгоритм со следующими характеристиками:

1. Использует SHA-256 как криптографическую хеш-функцию
2. Проверяет наличие определенного количества нулевых битов в начале хеша
3. Сложность настраивается через количество требуемых нулевых битов

Преимущества выбранного алгоритма:
- Простота реализации и верификации
- Настраиваемая сложность
- Асимметричность (решение требует значительно больше времени, чем проверка)
- Отсутствие зависимости от внешних ресурсов
- Детерминированная проверка

Дополнительные меры защиты:
- Защита от replay-атак через хранение использованных решений
- Ограничение времени жизни решений
- Rate limiting по IP
- Ограничение количества одновременных соединений

Для запуска выполните make up