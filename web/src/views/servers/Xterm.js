import { Terminal } from 'xterm'
import * as fit from 'xterm-addon-fit'
import * as attach from 'xterm-addon-attach'

Terminal.loadAddon(fit)
Terminal.loadAddon(attach)

export default Terminal

