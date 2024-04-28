var MyVarParams = Vue.extend({
    template: `<el-main style="max-width: 940px;margin: 0 auto 40px;">
        <el-row>
            <h3>模板变量的作用</h3>
            <p>在流水线使用的过程中，经常会出现创建重复的组件任务，而唯一不同的就是其中的参数名称。为了提高任务组件的复用性特引入了模板语法。</p>
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
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 根据 if 条件判断是否输出某个内容</p>
            <p><code>[[if .Condition]]...[[end]]</code></p>
            <p> &nbsp; &nbsp; &nbsp; &nbsp; 根据 if 条件判断是否输出某个内容</p>
            <p>官方内置函数:
                <ul>
                    <li><b>len</b> &nbsp; 返回一个字符串、数组、切片、映射或通道的长度。</li>
                    <li><b>index</b> &nbsp; 返回一个字符串、数组或切片中指定位置的元素。</li>
                    <li><b>printf</b> &nbsp; 根据格式字符串输出格式化的字符串。</li>
                    <li><b>range</b> &nbsp; 遍历一个数组、切片、映射或通道，并输出其中的每个元素。</li>
                    <li><b>with</b> &nbsp; 设置当前上下文中的变量。</li>
                </ul>
            </p>
            
            <p>应用内置函数：
                <ul>
                    <li><b>jsonString</b> &nbsp; json数据转字符串，示例：{"a":"A","b":{"b1":"B1","B2":22},"c":["c1","c2"]}</li>
                    <li><b>jsonString2</b> &nbsp; json数据转字符串2次，示例: {\\"a\\":\\"A\\",\\"b\\":{\\"b1\\":\\"B1\\",\\"B2\\":22},\\"c\\":[\\"c1\\",\\"c2\\"]}</li>
                </ul>
            </p>
            
        </el-row>
        
        <el-row>
            <h3>案例</h3>
            <p>下面任务设置了参数 <code>a</code>、<code>b</code> 并对参数进行了文字说明。</p>
            <p>在cmd类型任务下，进行了参数的具体使用；分别展示了常规参数使用、jsonString2函数对参数进行处理、b.b1二级参数的使用。</p>
            <p><i class="el-icon-warning"></i> 只有申明过出参数，才能在任务调用；</p>
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
            code:"echo A '\\\\n' map[B2:22 b1:B1] {\\\"B2\\\":22,\\\"b1\\\":\\\"B1\\\"} B1"
        }
    },
    // 模块初始化
    created(){
    },
    // 模块初始化
    mounted(){

    },

    // 具体方法
    methods:{

    }
})

Vue.component("MyVarParams", MyVarParams);