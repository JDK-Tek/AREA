import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

const port = 1234;

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  console.log("=> server listens on port " + port);
  await app.listen(1234);
}
bootstrap();
