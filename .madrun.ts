export default {
    'test': () => 'task test',
    'coverage': () => 'task coverage',
    'build': () => 'task build',
    'lint': () => 'putout .',
    'fix:lint': () => 'putout . --fix && task fix:lint',
};
