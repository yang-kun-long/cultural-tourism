const fs = require('fs');
const path = require('path');

// =================配置区域=================
// 扫描的根目录
const ROOT_DIR = process.cwd();

// 需要忽略的目录或文件 (支持正则)
const IGNORE_PATTERNS = [
    /^\./,              // 忽略以 . 开头的文件/目录 (如 .git, .idea, .env)
    /^node_modules$/,   // 忽略 node_modules
    /^tmp$/,            // 忽略 air 产生的临时目录
    /^docs$/,           // 忽略 swagger 生成的大量文档(可视情况保留)
    /^tests$/,          // 忽略测试目录
    /\.exe$/,           // 忽略编译后的可执行文件
    /\.log$/,           // 忽略日志
    /^go\.sum$/,        // 忽略依赖校验文件
    /^LICENSE$/,
    /^README\.md$/,
    /^scan_project\.js$/ // 忽略脚本自身
];

// 允许显示的特定 . 开头文件 (白名单)
const WHITELIST = [
    '.env.example',
    '.gitignore'
];
// =========================================

function shouldIgnore(name) {
    if (WHITELIST.includes(name)) return false;
    return IGNORE_PATTERNS.some(pattern => pattern.test(name));
}

function scanDir(dir, prefix = '') {
    let output = '';
    let entries;

    try {
        entries = fs.readdirSync(dir, { withFileTypes: true });
    } catch (err) {
        return '';
    }

    // 排序：文件夹在前，文件在后；同类按字母序
    entries.sort((a, b) => {
        if (a.isDirectory() && !b.isDirectory()) return -1;
        if (!a.isDirectory() && b.isDirectory()) return 1;
        return a.name.localeCompare(b.name);
    });

    // 过滤掉忽略的文件
    const filteredEntries = entries.filter(entry => !shouldIgnore(entry.name));

    filteredEntries.forEach((entry, index) => {
        const isLast = index === filteredEntries.length - 1;
        const marker = isLast ? '└── ' : '├── ';
        const childPrefix = isLast ? '    ' : '│   ';

        output += `${prefix}${marker}${entry.name}`;
        
        // 如果是文件，可以加一些备注（可选）
        if (entry.name === 'main.go') output += '  # [入口]';
        if (entry.name === 'client.go') output += '  # [核心] TCB SDK';
        
        output += '\n';

        if (entry.isDirectory()) {
            output += scanDir(path.join(dir, entry.name), prefix + childPrefix);
        }
    });

    return output;
}

console.log('Generating Project Structure...\n');
console.log('/');
console.log(scanDir(ROOT_DIR));