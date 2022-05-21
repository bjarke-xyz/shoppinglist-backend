import { listen, Router } from "worktop";

const router = new Router();

router.add("GET", "/", async (req, res) => {
  res.send(200, "yo");
});

router.add("GET", "/:name", async (req, res) => {
  const { name } = req.params;
  res.send(200, `yo 22${name}`);
});

listen(router.run);
