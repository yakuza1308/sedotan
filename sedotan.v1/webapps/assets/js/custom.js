function formatNumber(numb, format) {
    var ret = 0;
    // if(parseFloat(numb) >= Math.pow(10, 9)) {
    //     ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 9)), format) + ' B'; 
    // }
    // else if(parseFloat(numb) >= Math.pow(10, 6)) {
    //     ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 6)), format) + ' M'; 
    // }
    // else if(parseFloat(numb) >= Math.pow(10, 3)) {
    //     ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 3)), format) + ' K'; 
    // }
    // else {
    //     ret = kendo.toString(parseFloat(numb), format);    
    // }
     ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 6)), format); 

    return ret;
}

function NormalizeNumber(numb, format) {
    var ret = 0;
    if(parseFloat(numb) >= Math.pow(10, 9)) {
        ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 9)), format) + ' B'; 
    }
    else if(parseFloat(numb) >= Math.pow(10, 6)) {
        ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 6)), format) + ' M'; 
    }
    else if(parseFloat(numb) >= Math.pow(10, 3)) {
        ret = kendo.toString(parseFloat(numb)/(Math.pow(10, 3)), format) + ' K'; 
    }
    else {
        ret = kendo.toString(parseFloat(numb), format);    
    }

    return ret;
}

function GenDuration(startTime) {
    var ret = 0;

    var oneDay = 24*60*60*1000; // hours*minutes*seconds*milliseconds
    var firstDate = new Date();
    var secondDate = new Date(startTime);

    var diffDays = Math.round(Math.abs((firstDate.getTime() - secondDate.getTime())/(oneDay)));

    return diffDays;
}

function toUTC(d){
    var year = d.getFullYear();
    var month = d.getMonth();
    var date = d.getDate();
    var hours = d.getHours();
    var minutes = d.getMinutes();
    var seconds = d.getSeconds();
    return moment(Date.UTC(year, month, date, hours, minutes, seconds)).toISOString();
}
function jsonDate(strDt) {
    if (strDt == undefined) return "";
    var dt = str2date(strDt);
    if (dt.getFullYear() <= 1970 || dt.getFullYear() == 1) dt="";
    return dt;
}

function getUTCDate(strdate){
    var d = moment.utc(strdate);
    return new Date(d.year(), d.month(), d.date(), 0, 0, 0)
}
function toUTC(d){
    var year = d.getFullYear();
    var month = d.getMonth();
    var date = d.getDate();
    var hours = d.getHours();
    var minutes = d.getMinutes();
    var seconds = d.getSeconds();
    return moment(Date.UTC(year, month, date, hours, minutes, seconds)).toISOString();
}
function jsonDate(strDt) {
    if (strDt == undefined) return "";
    var dt = str2date(strDt);
    if (dt.getFullYear() <= 1970 || dt.getFullYear() == 1) dt="";
    return dt;
}

function jsonDateStr(dtSource, format) {
    var dt = str2date(dtSource);
    if (dt.getFullYear() <= 1970 || dt.getFullYear()==1) return "";
    return kendo.toString(dt, format==undefined ? jsonDateFormat : format);
}

function str2date(dtSource)
{
    dtSource = dtSource.toString();
    var dt = dtSource;
    if (dtSource.substr(0, 6) == "/Date(") {
        var dtParse = Date.parse(dtSource);
        if (isNaN(dtParse)) {
            var intMs = parseInt(dtSource.substr(6));
            dt = new Date(intMs);
        }
        else {
            dt = new Date(dtParse);
        }
        //alert(dt);
        dt = new Date(dt.getTime() + dt.getTimezoneOffset() * 60000);
    }
    else if(dtSource.length==5 && dtSource.substr(2,1)==":")
    {
        var times = dtSource.split(":");
        dt = new Date();
        dt = new Date(dt.getFullYear(), dt.getMonth(), dt.getDate(), times[0], times[1]);
    }
    else
    {
        dt = new Date(dtSource);
        dt = new Date(dt.getTime() + dt.getTimezoneOffset() * 60000);
    }
    return dt;
}

function date2str(dt, format) {
    if (dt instanceof Date === false) return "";
    if (dt.getFullYear() <= 1970 || dt.getFullYear() == 1) return "";
    return kendo.toString(dt, format == undefined ? jsonDateFormat : format);
}
