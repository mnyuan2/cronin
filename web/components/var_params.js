var MyVarParams = Vue.extend({
    template: `<el-main style="max-width: 940px;margin: 0 auto 40px;">
        <el-row>
            <h3>模板变量的作用</h3>
            <p>在流水线使用的过程中，经常会出现创建重复的组件任务，而唯一不同的就是其中的参数名称。为了提高任务组件的复用性特引入了模板语法。</p>
            <p>Golang <a target="_blank" href="https://juejin.cn/post/6844903762901860360">text/template</a></p>
        </el-row>
        
        <el-row>
            <h3>使用语法</h3>
            <p>在任务中定义了参数字段后，即可在任务详情中使用；变量由双中括号 [[ 和 ]] 包围，例如 [[.Name]] 表示一个名为 Name 的变量。变量可以是结构体、数组、切片、映射等类型，也可以是自定义类型。</p>
            <p><code>[[.Name]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 输出 Name 变量值</p>
            <p><code>[[len .Names]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp;使用函数获取 Names 变量的长度，并输出它的值</p>
            <p><code>[[range .Names]][[.]][[end]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 使用 range 语句遍历一个变量的值</p>
            <p><code>[[if .Condition]]...[[end]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 根据 if 条件判断是否输出某个内容，相关比较运算符：eq（==）、ne（!=）、lt（<）、le（<=）、gt（>）、ge（>=）</p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 完整样例：<code>[[if .Condition]] A [[else if .Condition]] B [[else]] C [[end]]</code></p>
            <p><code>[[$var := \`data\`]][[$var]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 模板内申明变量、使用变量;</p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 案例：<code>[[$data := json_decode \`{"code":"2","msg":"业务处理失败"}\`]][[if ne $data.code \`0\`]][[$data.msg]][[end]]</code> 输出 <code>业务处理失败</code></p>
            
            <h4>官方内置函数</h4>
            <p>
                <ul>
                    <li><b>len</b> &nbsp; 返回一个字符串、数组、切片、映射或通道的长度。</li>
                    <li><b>index</b> &nbsp; 返回一个字符串、数组或切片中指定位置的元素。</li>
                    <li><b>printf</b> &nbsp; 根据格式字符串输出格式化的字符串。
                        <ul><li>示例：<code>[[printf \`%v\` "666"]]</code> 输出 <code>666</code></li></ul>
                    </li>
                    <li><b>range</b> &nbsp; 遍历一个数组、切片、映射或通道，并输出其中的每个元素。</li>
                    <li><b>with</b> &nbsp; 设置当前上下文中的变量。</li>
                </ul>
            </p>
            
            <h4>系统内置函数</h4>
            <p>
                <ul>
<!--                    <li><b>jsonString</b> &nbsp; json数据转字符串，示例：{"a":"A","b":{"b1":"B1","B2":22},"c":["c1","c2"]}</li>-->
<!--                    <li><b>jsonString2</b> &nbsp; json数据转字符串2次，示例: {\\"a\\":\\"A\\",\\"b\\":{\\"b1\\":\\"B1\\",\\"B2\\":22},\\"c\\":[\\"c1\\",\\"c2\\"]}</li>-->
                    <li><b>json_decode<code>(string str) any</code></b> &nbsp; json字段解析，非规范json字符串抛出异常</li>
                    <li><b>rawurlencode<code>(string str) string</code></b> &nbsp; 字符串进行 url（RFC1738）编码，用于保护文字字符不被解释为特殊的URL分隔符。
                        <ul>
                            <li>示例：<code>[[rawurlencode \`Hello World!\`]]</code> 输出 <code>Hello%20World%21</code> 这里的空格被编码为%20，感叹号被编码为%21</li>
                        </ul>
                    </li>
                    <li><b>date<code>(string format=\`YYYY-MM-DD hh:mm:ss\` Time time=null) string</code></b> &nbsp; 格式化时间，参数为可选，默认为当前时间
                        <ul>
                            <li>参数1：format 输出格式表达式，字符串类型；默认值 YYYY-MM-DD hh:mm:ss；更多表达式说明如下：<pre>{{date_format}}</pre></li>
                            <li>参数2：time 时间，Time 对象；默认为当前时间</li>
                            <li>返回值：格式化后的字符串</li>
                            <li>示例：<code>[[date]]</code> 输出 <code>{{date_time}}</code></li>
                            <li>示例：<code>[[date \`YYYY-MM-DD\`]]</code> 输出 <code>{{date_only}}</code></li>
<!--                            <li>示例：<code>[[date \`hh:mm:ss\` Time]]</code> 输出 <code>{{time_only}}</code></li>-->
                        </ul>
                    </li>
                    <li><b>time<code>(string duration=\`\`) Time</code></b> 获取时间对象
                        <ul>
                            <li>参数1：duration 相对持续时间，示例："300ms"、"-1.5h"、"2h45m". 有效的时间单位是 "ns", "us" (or "µs"), "ms", "s", "m", "h".</li>
                            <li>返回值：Time 对象</li>
                            <li>示例：<code>[[time]]</code> 输出 <code>{{date_time}}.2306929 +0800 CST m=+0.047068901</code></li>
                            <li>示例：<code>[[time \`-720h\`]]</code> 输出 <code>{{date_time_30day}}.2306929 +0800 CST m=+0.047068901</code></li>
                            <li>与 <b>date</b> 函数组合使用：<code>[[date \`\` (time \`-720h\`)]]</code> 输出 <code>{{date_time_30day}}</code></li>
                        </ul>
                    </li>
                    <li><b>errorf<code>(str string, args ...any) error</code></b> 错误描述，执行后模板解析将终止并返回指定错误信息。
                        <ul>
                            <li>示例：<code>[[errorf \`错误描述\`]]</code> 输出 <code>错误描述</code></li>
                            <li>示例：<code>[[errorf \`错误 %s %v\` \`数值\` 5]]</code> 输出 <code>错误 数值 5</code></li>
                        </ul>
                    </li>
                    <li><b>float64<code>(any val) float64</code></b> 转浮点数</li>
                    <li><b>string<code>(any val) string</code></b> 转字符串</li>
                    <li><b>str_replace_calc<code>(str string, regex string, expr string) string</code></b> 正则匹配字符串计算后替换
                        <ul>
                            <li>参数1：str 被搜索的值</li>
                            <li>参数2：regex 数值匹配表达式，采用最左匹配原则</li>
                            <li>参数3：expr 计算公式</li>
                            <li>示例：<code>[[str_replace_calc \`release_v3.5.87.25 abc\` \`(\\d+)(\\D*$)\` \`+1\`]]</code> 输出 <code>release_v3.5.87.26 abc</code> 方法找到了字符串最左的数值并进行+1运算</li>
                        </ul>   
                    </li>
                </ul>
            </p>
<!--            <p>系统内置变量：-->
<!--                <ul>-->
<!--                    <li><b>date_only</b> string 日期 YYYY-MM-DD，示例：{{date_only}}</li>-->
<!--                    <li><b>time_only</b> string 时间 HH:mm:ss，示例：{{time_only}}</li>-->
<!--                    <li><b>date_time</b> string 日期时间 YYYY-MM-DD HH:mm:ss，示例：{{date_time}}</li>-->
<!--                </ul>-->
<!--            </p>-->
            
        </el-row>
        
        <el-row>
            <h3>案例</h3>
            <p>下面任务设置了参数 <code>a</code>、<code>b</code> 并对参数进行了文字说明。</p>
            <p>在cmd类型任务下，进行了参数的具体使用；分别展示了常规参数使用、jsonString2函数对参数进行处理、b.b1二级参数的使用。</p>
            <p><i class="el-icon-warning"></i> 只有申明过参数才能在任务内使用（系统变量除外）；</p>
            <el-image src="/static/image/var_field.png" style="width: 800px;"></el-image>
            <br>
            <p>流水线任务中可实现参数的具体值，包含的任务中只要申明了对应参数名称就会被传入。</p>
            <p>多余的参数实现会被忽略，任务中为实现的参数默认为空字符串 </p>
            <el-image src="/static/image/var_param.png" style="width: 800px;"></el-image>
            <br>
            <p>上面样例中，cmd任务文本最终执行语句为</p>
            <el-alert :title="code" type="info" :closable="false" style="background-color: #f4f4f5;"></el-alert>
</el-row>
    </el-main>`,

    name: "MyVarParams",
    data(){
        return {
            form:{},
            code:"echo A '\\\\n' map[B2:22 b1:B1] {\\\"B2\\\":22,\\\"b1\\\":\\\"B1\\\"} B1",
            cur_date: null,
            date_only: "0000-00-00",
            time_only: "00:00:00",
            date_time: "0000-00-00 00:00:00",
            date_time_30day: "0000-00-00 00:00:00",
            date_format: "M    - month (1)\nMM   - month (01)\nMMM  - month (Jan)\nMMMM - month (January)\nD    - day (2)\nDD   - day (02)\nDDD  - day (Mon)\nDDDD - day (Monday)\nYY   - year (06)\nYYYY - year (2006)\nhh   - hours (15)\nmm   - minutes (04)\nss   - seconds (05)"
        }
    },
    // 模块初始化
    created(){
    },
    // 模块初始化
    mounted(){
        this.currentDate()
    },

    // 具体方法
    methods:{
        currentDate(){
            this.cur_date = new Date()
            this.date_only = this.cur_date.getFullYear()+'-'+this.cur_date.getMonth().toString().padStart(2,'0')+'-'+this.cur_date.getDate() .toString().padStart(2,'0')
            this.time_only = this.cur_date.getHours().toString().padStart(2,'0')+':'+this.cur_date.getMinutes().toString().padStart(2,'0')+':'+this.cur_date.getSeconds().toString().padStart(2,'0')
            this.date_time = this.date_only +' '+ this.time_only
            let date30 = new Date((new Date()).getTime() - (30 * 24 * 60 * 60 * 1000))
            this.date_time_30day = date30.getFullYear()+'-'+date30.getMonth().toString().padStart(2,'0')+'-'+date30.getDate().toString().padStart(2,'0') + ' ' +
                date30.getHours().toString().padStart(2,'0')+':'+date30.getMinutes().toString().padStart(2,'0')+':'+date30.getSeconds().toString().padStart(2,'0')
            console.log(this.cur_date.toISOString(), date30.toISOString(), this.date_time_30day)
        }
    }
})

Vue.component("MyVarParams", MyVarParams);