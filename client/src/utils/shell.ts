import { exec } from "child_process";

export default class Shell {
    static execute(command: string, args: string[]) {
        return new Promise<string>((resolve, reject) => {
            exec(command + " " + args.join(" "), (error, stdout, stderr) => {
                if (error) {
                    reject(error.message)
                }
                if (stderr) {
                    reject(stderr)
                }
                if(stdout) {
                    resolve(stdout)
                }
            })
        })
    }
}