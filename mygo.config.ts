// Codie Bible: the Bible in many languages, its seasons and hymns, read aloud
// on a paired codie device.
export default ({ command }: { command: string }) => ({
  name: "Codie Bible",
  identifier: "io.github.codiebitcoin.codie-bibel",
  slug: "codie-bibel",
  version: "0.1.0",
  out: "build",
  urlSchemes: [command === "dev" ? "codie-bibel-dev" : "codie-bibel"],
  macos: {
    minimumSystemVersion: "13.0",
  },
  linux: {
    comment: "The Bible, its seasons and hymns, read aloud on your codie device",
    categories: ["Education"],
    webview: false,
  },
});
