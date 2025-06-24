module.exports = {
  apps: [{
    name: "muac-api",
    script: "./muac-api",
    watch: false,
    instances: 1,
    exec_mode: "fork",
    env: {
      NODE_ENV: "production",
      PORT: 8003
    }
  }]
}
