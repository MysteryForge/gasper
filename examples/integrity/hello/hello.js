import { sayHello } from "k6/x/gasper/integrity";

export default function () {
  const hello = sayHello();
  console.log(hello);
}
