FROM node:20-slim as base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base as install
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM base as build
COPY --from=install /app/node_modules ./node_modules
COPY . .
RUN pnpm run build

FROM base as production
COPY --from=install /app/node_modules ./node_modules
COPY --from=build /app/build ./build
COPY --from=build /app/.env ./
COPY --from=build /app/public ./public
COPY --from=build /app/package.json ./
COPY --from=build /app/server_ip.sh /scripts/server_ip.sh

USER node
EXPOSE 3000
CMD ["pnpm", "start"]