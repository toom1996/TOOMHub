Page({
  onShareAppMessage() {
    return {
      title: 'swiper',
      path: 'page/component/pages/swiper/swiper'
    }
  },

  data: {
    background: ['demo-text-1', 'demo-text-2', 'demo-text-3'],
    indicatorDots: true,
    vertical: false,
    autoplay: false,
    interval: 2000,
    duration: 200,
    image:[
      'http://toomhub.image.23cm.cn/tmp_da47d19c2206f05435765ff643218e4f6979f1f8d9858416.jpg',
      'http://toomhub.image.23cm.cn/tmp_596809151004b5e900823f16fd7f8772928afb16a2a3c273.jpg',
      'http://toomhub.image.23cm.cn/tmp_96867386fd2b66f0857172a5d8c113f72d37f824d8ba9d6f.jpg',
      'http://toomhub.image.23cm.cn/tmp_ee1f0d7cde40b9dcc48877ff05b8065387011e32b697a747.jpg',
      'http://toomhub.image.23cm.cn/tmp_b219880f212bb6aca099304f84a5a39724dbad542b2144ea.jpg'
    ]
  },

  changeIndicatorDots() {
    this.setData({
      indicatorDots: !this.data.indicatorDots
    })
  },

  changeAutoplay() {
    this.setData({
      autoplay: !this.data.autoplay
    })
  },

  intervalChange(e) {
    this.setData({
      interval: e.detail.value
    })
  },

  durationChange(e) {
    this.setData({
      duration: e.detail.value
    })
  },
  imageClickHandle() {
    console.log('sdfsdf')
  }
})