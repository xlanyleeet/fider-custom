# Gradient Fider

Наш fork: https://github.com/xlanyleeet/fider-custom
Оригінал: https://github.com/getfider/fider

`main` відстежує оригінальний Fider, `gradient` містить наші зміни.
Оригінальна історія та ліцензія AGPL-3.0 збережені.
Production на `feedback.gradient-mc.com` оновлюється окремо після перевірки
змін; створення цього форку саме по собі не змінює встановлений застосунок.

## Локальний запуск у Windows

Потрібен Docker Desktop у режимі Linux containers.

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\start-preview.ps1
```

Скрипт створює `deploy/.env.preview` із випадковими секретами, збирає образ
`gradient-fider:local` з поточного коду і запускає окремі PostgreSQL Alpine
та MailHog. Fider відкривається на http://localhost:3001,
тестові листи — на http://localhost:8026.
Ця база незалежна від локального `fider-docker` та сервера.

Повторний запуск готового образу:

```powershell
.\scripts\start-preview.ps1 -SkipBuild
```

Лише збірка Linux amd64 образу для нашого сервера:

```powershell
.\scripts\build-gradient.ps1
```

Логи та зупинка:

```powershell
docker compose --env-file deploy/.env.preview -f deploy/compose.preview.yaml logs -f app
docker compose --env-file deploy/.env.preview -f deploy/compose.preview.yaml down
```

`down` зберігає базу, `down -v` видаляє дані локальної preview-бази.
Справжні SMTP-паролі, ключі та дампи бази не додавай до цього репозиторію.

## Власні зміни

```powershell
git switch gradient
git switch -c feature/my-change
# Внести зміни, перевірити preview та відповідні тести.
git add <files>
git commit -m "Describe the change"
git push -u origin feature/my-change
```

Pull request створюй у `xlanyleeet/fider-custom` з базовою гілкою `gradient`.
Frontend: `public/`; backend: `app/`; переклади: `locale/`.
Порядок розробки й тести описано в `CONTRIBUTING.md`.

## Оновлення з оригіналу

Починай із чистого робочого дерева:

```powershell
git fetch upstream --tags
git switch main
git merge --ff-only upstream/main
git push origin main
git switch gradient
git merge main
# Вирішити конфлікти, виконати перевірки, потім git push origin gradient.
```

Перед production-оновленням зроби резервну копію бази та перевір міграції.
На малому Oracle-сервері використовуй готовий образ, зібраний локально або в CI.
Налаштування чинного розгортання лежать окремо, у сусідній папці `fider-docker`.
