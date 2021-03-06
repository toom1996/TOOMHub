const app = getApp();
Page({

  /**
   * 页面的初始数据
   */
  data: {
    host: '',
    videoUrl: '',
    duration: 0,
    coverTotalFrame: 40, //封面总帧数
    coverInterval: 1, //封面间隔
    currentCover: '', //当前选择的封面
    checkedCover: 0, //默认选择第0帧
    coverHeight: ( app.globalData.windowWidth - 40 ) / 4 ,
    coverType: 0 //选择封面的类型
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    console.log(options.host + options.name + '?vframe/jpg/offset/' + this.data.checkedCover)
    this.setData({
      videoUrl: options.host + options.name,
      duration: parseInt(options.duration),
      defaultCover: options.host + options.name + '?vframe/jpg/offset/' + this.data.checkedCover,
    })
    this.setData({
      coverInterval: this.data.duration / 40
    })
    console.log(this.data.coverInterval)
    console.log(this.data.duration)
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  },
  selectCoverHandle(event) {
    console.log(event)
    let index = event.currentTarget.dataset.index;
    this.setData({
      checkedCover: index,
    })
    this.setData({
      defaultCover: this.data.videoUrl + '?vframe/jpg/offset/' + this.data.checkedCover * this.data.coverInterval
    })
  },
  
  //选择封面处理事件
  checkedCoverHandel() {
    console.log(this.data.coverType)
    if (this.data.coverType == 1) {
      console.log(app.getApi('requestHost') + '/c/image/upload')
      wx.uploadFile({
        url: app.getApi('requestHost') + '/c/image/upload', //仅为示例，非真实的接口地址
        filePath: this.data.defaultCover,
        name: 'file',
        success (res){
          console.log(res)
          let resData = JSON.parse(res.data);
          app.globalData.videoCover = resData.data.request_host + resData.data.name;
          console.log('videoCover', app.globalData.videoCover)
          wx.navigateBack({
            delta: -1
          })
          //do something
        }
      })
      
    } else {
      app.globalData.videoCover = '?vframe/jpg/offset/' + this.data.checkedCover * this.data.coverInterval;
      wx.navigateBack({
        delta: -1
      })
    }
  },


  afterRead(event){
    console.log('afterRead')
    let file = event.detail.file.path;
    this.setData({
      defaultCover: file,
      coverType:1
    })
  }
})