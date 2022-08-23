/**
 * 时间戳 转 date
 * @param nS
 * @returns {string}
 */
function getLocalTime(nS) {
    // 10 需要 * 1000； 13位不做任何操作。
    return new Date(parseInt(nS)).toLocaleString().replace(/:\d{1,2}$/,' ');
}

/* 质朴长存法 by lifesinger */
function pad(num, n) {
    var len = num.toString().length;
    while(len < n) {
        num = "0" + num;
        len++;
    }
    return num;
}

// 通过 date 对象 获取date 时间字符串
function getDateString(d) {
    return d.getFullYear()+'-'+pad(d.getMonth()+1,2)+'-'+pad(d.getDate(), 2)
}

// 通过 date 对象 获取datetime 时间字符串
function getDatetimeString(d) {
    return d.getFullYear()+'-'+pad(d.getMonth()+1,2)+'-'+pad(d.getDate(), 2)+ ' '+ pad(d.getHours())+':'+d.getMinutes()+':'+d.getSeconds()
}

/**
 * 网络请求
 * @type {{post: api.post, get: api.get}}
 */
var api = {
    baseUrl : window.location.protocol+"//"+window.location.host,
    /**
     * 内网get请求
     * @param path string 请求路径
     * @param param object 请求参数 将拼url参数
     * @param success func 响应结果
     */
    innerGet: function (path,param, success) {
        let url =  this.baseUrl+path
        if (param){
            url += "?"
            for (i in param){
                if (param[i]){
                    url += i+"="+param[i]+"&"
                }
            }
            url = url.slice(0,-1)
        }
        console.log("get", url)
        this.ajax(url, null, null,null, success, null, true, 'get')
    },

    /**
     * 内网post请求
     */
    innerPost: function (path, param, success) {
        let url =  this.baseUrl+path
        param = JSON.stringify(param)
        console.log("post", url, param)
        this.ajax(url, param, null,null, success, null, true, 'post', 'json')
    },

    innerGetFile: function(path){
        let features = "height=240, width=400, top=50, left=50, toolbar=no, menubar=no,scrollbars=no,resizable=no, location=no, status=no,toolbar=no";
        window.open(this.baseUrl+path, "_blank", features)
    },

    ajax: function (url, data, beforeSend, complete, success, error, async=true, type='post', dataType='json') {
        var success = success || function (data) {
            // 默认Success方法
            if(data.code == '-403'){
                alert(res.msg);
                return ;
            }
        };

        var error = error || function (data) {
            // 默认Error方法
            console.log(data);
            setTimeout(function () {

                if (data.status == 404) {
                    alert('请求不存在404');
                } else if(data.status == 503) {
                    alert('请求失败503');
                } else if(data.status == 401){
                    alert('登录超时');
                    setTimeout(()=>{location.reload();},1000);
                }else if (data.message != ""){
                    alert(data.message);
                }else{
                    alert('请求失败,网络连接超时');
                }
                ajaxStatus = true;
            }, 500);
        };
        $.ajax({
            'url': url,
            'data': data,
            'type': type,
            'headers': {
                'Content-Type': 'application/json'
            },
            'dataType': dataType,
            'async': async,
            'beforeSend': beforeSend,
            'complete': complete,
            'success': success,
            'error': error
        });
    }
}