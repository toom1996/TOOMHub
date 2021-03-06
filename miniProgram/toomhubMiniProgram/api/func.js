const app = getApp()

// 生成缩略图
export const getThumbnail = (options, callback) => {
  wx.showLoading({
    title: '分享图片生成中',
  })
  let imgData = [];
  console.log('line 8 ---->', JSON.parse(JSON.stringify(imgData)))
  let dataset = options.target.dataset;
  console.log('dataset -->', dataset)
  let host = 'https://toomhub.23cm.cn/image/';
  let shareData = {
    title: dataset.title, //分享内容
    id: dataset.id, //分享的id
    type: dataset.type, //分享类型
    createdBy: dataset.createdby, //创建人
    avatar: dataset.avatar,
  }

  //video
  if (dataset.type == 1) {
    var internetImgData = [ host + dataset.video + dataset.cover, dataset.avatar]
  }

  //image 
  if (dataset.type == 0) {
    var internetImgData = [host + dataset.cover + '?imageView2/1/q/75/format/jpg', dataset.avatar]
    console.log('image internetImgData -> ',internetImgData)
  }

  console.log('line 26 ---->', JSON.parse(JSON.stringify(imgData)))
  for (let i = 0; i < internetImgData.length; i++) {
    //getImageInfo 不支持webp格式文件
    wx.getImageInfo({
      src: internetImgData[i],
      complete: (res) => {
        console.log('res ->', res)
        console.log(internetImgData[i])
        imgData[i] = {
          path: res.path,
          width: res.width,
          height: res.height,
        };
        if (imgData.hasOwnProperty(0) === true && imgData.hasOwnProperty(1) == true) {
          createPhoto(shareData, imgData, callback)
        }
      },fail: (e) => {
        wx.hideLoading({
          success: (res) => {},
        })
        wx.showToast({
          title: '图片分享失败啦~~',
        })
      }
    })
  }
  console.log('complete')
}


function createPhoto(data, imgData, callback) {
  try{
    let cavansId = 'shareCanvas';
  let createdBy = '匿名用户';
  //视频分享

  const context = wx.createCanvasContext(cavansId, self);
  //截取用户名
  if (strlen(data.createdBy) > 14) {
    createdBy = data.createdBy.substring(0, 14) + '...';
  }

  //----- 绘制头像 -----
  var avatarurl_width = 45;    //绘制的头像宽度
  var avatarurl_heigth = 45;   //绘制的头像高度
  var avatarurl_x = 0;   //绘制的头像在画布上的位置
  var avatarurl_y = 0;   //绘制的头像在画布上的位置
  context.save();
  context.beginPath(); //开始绘制
  context.arc(avatarurl_width / 2 + avatarurl_x, avatarurl_heigth / 2 + avatarurl_y, avatarurl_width / 2, 0, Math.PI * 2, false);
  context.clip()	//裁剪
  context.drawImage(imgData[1].path, avatarurl_x, avatarurl_y, avatarurl_width, avatarurl_heigth);
  context.restore(); //恢复之前保存的绘图上下文 恢复之前保存的绘图上下午即状态 还可以继续绘制

  //----- 绘制用户昵称 -----
  context.setFontSize(20)
  context.fillText(createdBy, 55, 30)
  // ----- 绘制封面图 -----
  console.log(imgData)
  let cWidth = 500; //画布宽度
  let cHeight = 400; //画布高度
  let imgWidth = imgData[0].width; //封面宽度
  let imgHeight = imgData[0].height; //封面高度
  let dWidth = cWidth / imgWidth;  // canvas与图片的宽度比例
  let dHeight = cHeight / imgHeight;  // canvas与图片的高度比例

  if (imgWidth > cWidth && imgHeight > cHeight || imgWidth < cWidth && imgHeight < cHeight) {
    if (dWidth > dHeight) {
      context.drawImage(imgData[0].path, 0, (imgHeight - cHeight / dWidth) / 2, imgWidth, cHeight / dWidth, 0, 50, cWidth, cHeight)
    } else {
      context.drawImage(imgData[0].path, (imgWidth - cWidth / dHeight) / 2, 50, cWidth / dHeight, imgHeight, 0, 50, cWidth, cHeight)
    }
  } else {
    if (imgWidth < cWidth) {
      context.drawImage(imgData[0].path, 0, (imgHeight - cHeight / dWidth) / 2, imgWidth, cHeight / dWidth, 0, 50, cWidth, cHeight)
    } else {
      context.drawImage(imgData[0].path, (imgWidth - cWidth / dHeight) / 2, 20, cWidth / dHeight, imgHeight, 0, 50, cWidth, cHeight)
    }
  }

  if (data.type == 1) {
    // ----- 绘制视频播放按钮 -----
    context.setGlobalAlpha(0.5);
    context.drawImage('/static/icon/play.png', 136, 111, 200, 200);
    context.setGlobalAlpha(1);
  }

  // ----- 绘制图片 -----
  context.draw(false, () => {
    wx.canvasToTempFilePath({
      canvasId: cavansId,
      x: 0,
      y: 0,
      width: cWidth,
      height: cHeight,
      complete: (tmpRes) => {
        callback({
          src: tmpRes.tempFilePath,
          id: data.id,
          title: data.title
        });
        wx.hideLoading({
        })
      }
    }, this);
  })
  }catch(e){
    wx.hideLoading({
      success: (res) => {
        wx.showToast({
          title: '图像生成错误',
        })
      },
    })
  }
}


export const strlen = (str) => {
  var len = 0;
  for (var i = 0; i < str.length; i++) {
    var c = str.charCodeAt(i);
    //单字节加1
    if ((c >= 0x0001 && c <= 0x007e) || (0xff60 <= c && c <= 0xff9f)) {
      len++;
    }
    else {
      len += 2;
    }
  }
  return len;
}

export const calculateVideoSize = (width, height, paddingPx = 32) => {
  width = parseInt(width)
  height = parseInt(height)
  let screenWidth = app.globalData.windowWidth;
  let screenHeight = app.globalData.windowHeight;
  // let paddingPx = 32; //左右边距
  let calculateWidth = 0;
  let calculateHeight = 0;
  app.globalData.windowHeight

  //宽大于屏幕宽度, 高大于最大高度/2, 则压缩高度 和宽度
  if (width >= (screenWidth - paddingPx) && height >= screenHeight / 2) {
    if (width < height) {
      console.log('mode 1 - 1')
      let p = height / (screenHeight / 2)
      console.log(p)
      calculateHeight = height / p
      calculateWidth = width / p
    } else {
      let p = width / (screenWidth - paddingPx)
      calculateHeight = height / p
      calculateWidth = width / p
    } 
  }

  //宽大于屏幕宽度, 高大于最大高度/2, 则压缩高度 和宽度
  if (width >= (screenWidth - paddingPx) && height <= screenHeight / 2) {
    console.log('mode 2')
    let p = width / (screenWidth - paddingPx)
    calculateHeight = height / p
    calculateWidth = width / p  
  }

  //宽大于屏幕宽度, 高大于最大高度/2, 则压缩高度 和宽度
  if (width <= (screenWidth - paddingPx) && height >= screenHeight / 2) {
    console.log('mode 3')
    let p = height / (screenHeight / 2)
    calculateHeight = height / p
    calculateWidth = width / p
  }

  //宽大于屏幕宽度, 高大于最大高度/2, 则压缩高度 和宽度
  if (width <= (screenWidth - paddingPx) && height <= screenHeight / 2) {
    console.log('mode4')
    let p = width / (screenWidth - paddingPx)
    calculateHeight = height / p
    calculateWidth = width / p
  }



  return {
    width: calculateWidth,
    height:calculateHeight
  }
}

