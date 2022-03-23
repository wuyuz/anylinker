<template>
  <div  style="height: 100%;"
       id="terminal"></div>
</template>
<script>
import { Terminal } from 'xterm';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { Message } from "element-ui";

export default {
  name: 'Console',
  props: {
    terminal: {
      type: Object,
      // eslint-disable-next-line vue/require-valid-default-prop
      default: {}
    },
    url: {
      type: String,
      default:""
    }
  },
  data () {
    return {
      term: null,
      terminalSocket: null,
      rows: 40,
      cols: 100
    }
  },
  methods: {
    runRealTerminal () {
      console.log('webSocket is finished')
    },
    errorRealTerminal () {
      console.log('error')
    },
    closeRealTerminal () {
      console.log('close')
    }
  },
  mounted () {
    this.openTer();
  },
  beforeDestroy () {
    this.terminalSocket.close();
    this.term=null;
  },
  methods:{
    openTer() {
      let terminalContainer = document.getElementById('terminal')
      const _this = this;
      this.term = new Terminal({
        rendererType: 'canvas', // 渲染类型
        rows: parseInt(_this.rows), // 行数
        cols: parseInt(_this.cols), // 不指定行数，自动回车后光标从下一行开始
        convertEol: true, // 启用时，光标将设置为下一行的开头
        scrollback: 100, // 终端中的回滚量
        disableStdin: false, // 是否应禁用输入。
        cursorStyle: 'underline', // 光标样式
        cursorBlink: true, // 光标闪烁
        theme: {
          foreground: '#7e9192', // 字体
          background: 'black', // 背景色
          cursor: 'help', // 设置光标
          lineHeight: 16
        },
        windowOptions:{
            fullscreenWin: true
        }
      })

      this.term.open(terminalContainer)
      // open websocket
      this.terminalSocket = new WebSocket(this.url)
      // this.terminalSocket.onopen = this.runRealTerminal
      // this.terminalSocket.onclose = this.closeRealTerminal
      // this.terminalSocket.onerror = this.errorRealTerminal

      this.terminalSocket.onerror = () => {
          Message.error('ws has no token, please login first')
      }

      this.terminalSocket.onclose = () => {
          this.term.setOption('cursorBlink', false)
          Message.warning('console.web_socket_disconnect')
      }

      const attachAddon = new AttachAddon(this.terminalSocket);
      this.term.loadAddon(attachAddon);
      const fitAddon = new FitAddon();
      fitAddon.fit();
      this.term._initialized = true
    }
  },
  watch:{
    url(v,old){
      if (v !== old) {
        // this.terminalSocket.onclose();
        this.term = null;
        if (v !== "") {
          this.openTer()
        }
      }
    }
  }
}
</script>
