console.log("SEEDING DATABASE ########################");

db = db.getSiblingDB(process.env.MONGO_DB_NAME);

db.createUser({
  user: process.env.MONGO_USERNAME,
  pwd: process.env.MONGO_PASSWORD,
  roles: [{ role: "readWrite", db: process.env.MONGO_DB_NAME }],
});

db.createCollection("users");
db.createCollection("doujins");

db.users.insertMany([
  {
    _id: ObjectId("66aa96e9f3fcf7602ce6ad68"),
    discord_id: BigInt("11111111111111111"),
    name: "InfernalHydra",
    reservations: [
      {
        doujin_id: ObjectId("66aa96e8f3fcf7602ce6ad66"),
        datetime_added: new Date("2024-07-31T19:56:25.845Z"),
      },
      {
        doujin_id: ObjectId("66aa96e9f3fcf7602ce6ad67"),
        datetime_added: new Date("2024-07-31T19:56:26.085Z"),
      },
    ],
    last_updated: new Date("2024-07-31T19:56:26.085Z"),
  },
]);

db.doujins.insertMany([
  {
    _id: ObjectId("66aa96e8f3fcf7602ce6ad66"),
    title: "【C104新作+過去作まとめ買い】Lapin「KAZAMARINE Collection」セット",
    price_in_yen: 3331,
    price_in_usd: 22.317700000000002,
    image_preview_url:
      "https://melonbooks.akamaized.net/user_data/packages/resize_image.php?image=218001043810.jpg",
    url: "https://www.melonbooks.co.jp/detail/detail.php?product_id=2479022",
    is_r18: false,
    circle_name: "Lapin",
    author_names: ["宮城らとな"],
    genres: ["ホロライブ", "バーチャルユーチューバー"],
    events: ["コミックマーケット104"],
    last_updated: new Date("2024-07-31T19:56:25.845Z"),
    reservations: [
      {
        user_id: ObjectId("66aa96e9f3fcf7602ce6ad68"),
        datetime_added: new Date("2024-07-31T19:56:25.845Z"),
      },
    ],
  },
  {
    _id: ObjectId("66aa96e9f3fcf7602ce6ad67"),
    title: "LACKGAKIBOX",
    price_in_yen: 3960,
    price_in_usd: 26.532,
    image_preview_url:
      "https://melonbooks.akamaized.net/user_data/packages/resize_image.php?image=212001450086.jpg",
    url: "https://www.melonbooks.co.jp/detail/detail.php?product_id=2481050",
    is_r18: false,
    circle_name: "珈琲紳士の部屋",
    author_names: ["lack"],
    genres: ["Fate/Grand Order", "刀剣乱舞"],
    events: ["コミックマーケット104"],
    last_updated: new Date("2024-07-31T19:56:26.085Z"),
    reservations: [
      {
        user_id: ObjectId("66aa96e9f3fcf7602ce6ad68"),
        datetime_added: new Date("2024-07-31T19:56:26.085Z"),
      },
    ],
  },
]);

console.log("SEEDING COMPLETE ########################");
