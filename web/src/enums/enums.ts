export enum Mode {
  ECH0 = 0, // 默认编辑状态
  Panel = 1, // 显示面板状态
  TODO = 2, // 待办事项状态
  EXTEN = 3, // 处理扩展状态
  PlayMusic = 4, // 音乐播放器状态
  Image = 5, // 图片上传状态
}

export enum ExtensionType {
  MUSIC = 'MUSIC',
  VIDEO = 'VIDEO',
  GITHUBPROJ = 'GITHUBPROJ',
  WEBSITE = 'WEBSITE',
}

export enum ImageSource {
  LOCAL = 'local',
  URL = 'url',
  S3 = 's3',
  R2 = 'r2',
}

export enum CommentProvider {
  TWIKOO = 'twikoo',
  ARTALK = 'artalk',
  WALINE = 'waline',
  GISCUS = 'giscus',
}
