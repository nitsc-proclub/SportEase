# Backend

## ディレクトリ構造
```
backend/
├── cmd/
│   └── server/
│       └── main.go        # エントリーポイント
├── config/                # 設定ファイル（DB, 環境変数など）
│   └── config.go
├── internal/
│   ├── handler/           # HTTPハンドラ (コントローラ)
│   │   └── user_handler.go
│   ├── service/           # ビジネスロジック
│   │   └── user_service.go
│   ├── repository/        # DBアクセス (SQL/ORM)
│   │   └── user_repository.go
│   ├── middleware/        # ミドルウェア
│   │   └── auth_middleware.go
│   ├── model/             # モデル（構造体定義）
│   │   └── user.go
│   └── router/            # ルーティング設定
│       └── router.go
├── pkg/                   # 再利用可能なユーティリティ (JWT, Logger, etc.)
│   └── logger/logger.go
├── go.mod
└── go.sum
```

## API設計

### 1. 認証

#### `POST /api/auth/login`
- **説明:** ユーザーを認証し、JWTを返します。
- **リクエストボディ:**
  ```json
  {
    "username": "ユーザー名または管理者名",
    "password": "パスワード"
  }
  ```
- **レスポンス:**
  - `200 OK`: ログイン成功
    ```json
    {
      "token": "あなたのJWTトークン"
    }
    ```
  - `401 Unauthorized`: 認証情報が無効です。

#### `POST /api/auth/logout`
- **説明:** ユーザーをログアウトします（クライアント側でトークンを削除することで実装可能）。
- **認証:** 必須
- **レスポンス:**
  - `200 OK`: ログアウト成功

#### `GET /api/auth/me`
- **説明:** 認証されたユーザーの情報を取得します。
- **認証:** 必須
- **レスポンス:**
  - `200 OK`: ユーザー情報

### 2. トーナメント

#### `GET /api/tournaments/:sport`
- **説明:** 指定された競技のトーナメント情報を取得します。天候によって結果を絞り込むことも可能です。
- **認証:** 全ユーザー（student, admin, superroot）で必須
- **URLパラメータ:**
  - `sport`: 競技名 (例: `volleyball`, `soccer`)
- **クエリパラメータ:**
  - `weather` (任意): `sunny` または `rainy`。指定された天候用のトーナメントを返します。指定がない場合、天候条件が `any` のものと、`sunny` (デフォルト) のものを結合して返します。
- **レスポンス:**
  - `200 OK`: 指定された競技・天候に合致するトーナメント情報の配列
    ```json
    [
      {
        "id": 1,
        "name": "令和7年度 秋季大会 バレーボール",
        "sport": "volleyball",
        "weather_condition": "sunny",
        "teams": [
          {"id": 1, "name": "3A"},
          {"id": 2, "name": "3B"}
        ],
        "matches": [
          {
            "id": 1,
            "round": 1,
            "match_number_in_round": 1,
            "team1_id": 1,
            "team2_id": 2,
            "team1_score": 0,
            "team2_score": 0,
            "winner_team_id": null,
            "next_match_id": 3
          }
        ]
      }
    ]
    ```
  - `404 Not Found`: 指定された競技のトーナメントが見つかりません。

### 3. 試合

#### `PUT /api/matches/:id`
- **説明:** 特定の試合のスコアを更新します。勝者はバックエンドで自動的に計算・更新されます。
- **認証:** 必須 (admin, superroot)
  - `admin` は担当する競技の試合のみ更新できます。
  - `superroot` はどの試合でも更新できます。
- **リクエストボディ:**
  ```json
  {
    "user": "更新者",
    "team1_score": 2,
    "team2_score": 1,
    "winner_team_id": 1,
    "status": "finished"
  }
  ```
- **レスポンス:**
  - `200 OK`: 更新成功。更新された試合オブジェクトを返します。
  - `403 Forbidden`: この試合を更新する権限がありません。
  - `404 Not Found`: 指定されたIDの試合は存在しません。

#### `GET /api/matches/:sport`
- **説明:** 指定された競技の試合情報を取得します。
- **認証:** admin, superroot
- **URLパラメータ:**
  - `sport`: 競技名 (例: `volleyball`, `soccer`)
- **レスポンス:**
  - `200 OK`: 指定された競技の試合情報の配列
    ```json
    [
      {
        "id": 1,
        "match_number_in_round": 1,
        "round": 1,
        "team1_id": 10,
        "team1_name": "A組",
        "team2_id": 11,
        "team2_name": "B組",
        "team1_score": 2,
        "team2_score": 1,
        "winner_team_id": 10,
        "status": "finished",
        "next_match_id": 5
      },
  ]
     - `404 Not Found`: 指定された競技の試合が見つかりません。

### 4. 管理者機能 (Superrootのみ)

#### `GET /api/admin/users`
- **説明:** すべての `admin` ユーザーのリストを取得します。
- **認証:** 必須 (superrootのみ)
- **レスポンス:**
  - `200 OK`:
    ```json
    [
      {
        "id": 2,
        "username": "admin_volleyball",
        "role": "admin",
        "assigned_sport": "volleyball"
      }
    ]
    ```
  - `403 Forbidden`: superrootユーザーではありません。

#### `POST /api/admin/users`
- **説明:** 新しい `admin` ユーザーを作成します。
- **認証:** 必須 (superrootのみ)
- **リクエストボディ:**
  ```json
  {
    "username": "new_admin",
    "password": "new_password",
    "assigned_sport": "soccer"
  }
  ```
- **レスポンス:**
  - `201 Created`: ユーザーの作成に成功しました。
  - `403 Forbidden`: superrootユーザーではありません。
  - `400 Bad Request`: 無効な入力です（例: フィールドの欠落）。

#### `PUT /api/admin/users/:id`
- **説明:** `admin` ユーザーの情報を更新します。
- **認証:** 必須 (superrootのみ)
- **リクエストボディ:**
  ```json
  {
    "username": "updated_admin_name",
    "password": "optional_new_password",
    "assigned_sport": "table_tennis"
  }
  ```
- **レスポンス:**
  - `200 OK`: ユーザーの更新に成功しました。
  - `403 Forbidden`: superrootユーザーではありません。
  - `404 Not Found`: 指定されたIDのユーザーは存在しません。

#### `DELETE /api/admin/users/:id`
- **説明:** `admin` ユーザーを削除します。
- **認証:** 必須 (superrootのみ)
- **レスポンス:**
  - `204 No Content`: ユーザーの削除に成功しました。
  - `403 Forbidden`: superrootユーザーではありません。
  - `404 Not Found`: 指定されたIDのユーザーは存在しません。
