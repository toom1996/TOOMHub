// pages/components/tabbar/tabbar.js
Component({
  /**
   * 组件的属性列表
   */
  properties: {
  },

  /**
   * 组件的初始数据
   */
  data: {
    active: 'pages/index/index'
  },
  // 组件生命周期：一打开页面就执行
  attached:function (){
    let pages = getCurrentPages()
    let currentPages = pages[pages.length - 1]
    this.setData({
      active: currentPages.route
    })
  },
  /**
   * 组件的方法列表
   */
  methods: {
    navigationSwitch:function (event) {
      wx.reLaunch({
        url: '/' + event.detail,
      })
    }
  }
})
