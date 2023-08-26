package boards

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/soil-demo/util"
)

/*POST[pt]
A ideia desse Automata é simular a dispersão de um
liquido em um meio solido permeável, como solo, areia, etc.
Ele foi criado como protótipo de um sistema de dispersão de
água subterrânea para meu jogo.

Como vai perceber, não existe muita complexidade nesse código,
essa é a magica desse tipo de algoritmo, com regras simples é
possível construir comportamentos complexos. Aqui, simulamos
a interação entre líquidos e sólidos no tempo e espaço usando apenas
aritmética básica.

Nosso espaço será representado como um 2d array de vetores, mas
no nosso caso eles servem apenas como uma forma conveniente de
agrupar nossos dois valores relevantes:

- Umidade
	- A quantidade de liquido contido naquela célula
- Impermeabilidade
	- A resistência que o material daquela célula apresenta ao movimento de liquido

Os campos `Rocks` e `Rain` servem para demonstrar, respectivamente,
células totalmente impermeáveis e células que sao fontes de umidade.
*/

/*POST[es]
The goal for this Automata is to simulate the diffusion of
a liquid on a solid medium, like soil, sand, etc. It was created
as a prototype of a underground water diffusion for my game.

As will notice, there isn't much complexity in this code, that's the magic
of this kind of algorithm, we can use simple rules to create complex
behavior. Here, we simulate the interaction between liquids ans solids
over time and space using basic arithmetic

Our board will be represented by a 2d array of vectors, this is just
for convenience since we can easily group the two relevant values:

- Humidity
	- The amount of liquid within the cell
- Impermeability
	- The resistance of the cell to the movement of liquid trough it

The fields `Rocks` and `Rain` just represent cells that are totally impermeable
or sources of humidity, respectively.
*/

// PIN
type HumidityBoard struct {
	initValues  [][]mgl32.Vec2 // [humidity, impermeability]
	values      [][]mgl32.Vec2 // [humidity, impermeability]
	Rocks, Rain [][]bool
	hvrX, hvrY  int
}

func (ba *HumidityBoard) Size() (int, int) {
	return len(ba.values), len(ba.values[0])
}

/*POST[pt]
A lógica para determinar a umidade de uma célula é bastante simples,
entendendo que ao longo do tempo o liquido tende a se espalhar uniformemente
pelo espaço, assumimos que o nosso valor de umidade da célula é a média aritmética
entre o seu valor atual e a das suas vizinhas. Por simplicidade assumimos apenas
4 vizinhos conforme esse diagrama onde `v0` é a célula atual:

```
   |v3|
|v1|v0|v2|
   |v4|
```

Para considerar a permeabilidade, utilizamos uma média ponderada onde
o peso de cada termo é definido pelo valor de impermeabilidade daquela célula.
O detalhe é que consideramos o reciproco(1/valor) como peso das células vizinhas,
dessa forma conseguimos o efeito de que uma célula altamente impermeável tende
a não perder umidade ao mesmo tempo que resiste à absorção de mais liquido. Além
disso, verificamos se a célula excede o nosso limite de 1023 e ajustamos então o valor para esse limite.
*/

/*POST[es]
The logic to determine the humidity of a cell is simples, since the liquid
converges to a uniform distribution over time for all the space, we assume
that the new value of humidity of a cell is just the mean between it and
their neighbors values. For simplicity, we only count 4 neighbors as the
diagram below, where `v0` is the current cell:

```
   |v3|
|v1|v0|v2|
   |v4|
```

To include permeability properties, we use a weighted mean where each
factor is defined by that cell impermeability. The catch is that we use
the inverse(1/value) for neighboring cells, that way we achieve the effect
of a highly impermeable cell resisting either losing or gaining humidity.
Then we clamp to new value to a max of 1023.
*/

// PIN
func (ba *HumidityBoard) Update() error {
	// Armazenamos o estado inicial do espaço para servir de referencia
	m0 := [][]mgl32.Vec2{}

	for x, row := range ba.values {
		m0 = append(m0, []mgl32.Vec2{})
		for y, v0 := range row {
			m0[x] = append(m0[x], v0)
			// Cell names
			//    |v3|
			// |v1|v0|v2|
			//    |v4|
			// Pulamos o calculo de fontes de umidade e células com alto impermeabilidade
			if ba.Rain[x][y] || v0[1] >= (math.MaxFloat32/5)*4 {
				continue
			}

			// Assumimos que as bordas são células secas e impermeáveis
			v1 := mgl32.Vec2{0, math.MaxFloat32}
			v2 := mgl32.Vec2{0, math.MaxFloat32}
			v3 := mgl32.Vec2{0, math.MaxFloat32}
			v4 := mgl32.Vec2{0, math.MaxFloat32}

			// Verificamos se o vizinho existe e aplicamos os valores corretos
			if x > 0 {
				v1 = ba.values[x-1][y]
			}
			if x < len(ba.values)-1 {
				v2 = ba.values[x+1][y]
			}
			if y > 0 {
				v3 = ba.values[x][y-1]
			}
			if y < len(ba.values[0])-1 {
				v4 = ba.values[x][y+1]
			}

			// Calculamos a média aritmética ponderada
			r := ((v0[0] * (v0[1])) + (v1[0] / v1[1]) + (v2[0] / v2[1]) + (v3[0] / v3[1]) + (v4[0] / v4[1])) / (v0[1] + (1 / v1[1]) + (1 / v2[1]) + (1 / v3[1]) + (1 / v4[1]))

			// Limitamos os valores a um máximo de 1023
			if r > 1023 {
				r = 1023
			}

			// Atualizamos espaço com novo valor de umidade
			m0[x][y][0] = r
		}
	}
	ba.values = m0
	return nil
}

/*POST[pt]
O resultado no fim não é perfeito, há parâmetros que não são levados em consideração como
velocidade e densidade do liquido, mas para o nosso caso já é suficiente. Outra limitação
é quanto a conservação de massa do sistema. Aos poucos o volume total de umidade cai e isso
causa o efeito de umidade desaparecendo espontaneamente, que é fisicamente impossível.

Como dito, esse é um protótipo e limitações como essa não são necessariamente problemas
para aplicação em jogos. Vale lembrar que esse algoritmo foi escrito de forma síncrona,
mas é totalmente possível adapta-lo para operar de forma paralelizada.
*/

/*POST[es]
The final result is not perfect, there are parameters that are not accounted like velocity
and density of the liquid, but it's sufficient for our porpoises. Another limitation of this
method is lack of conservation of mass. As the simulation evolves the humidity spontaneously
drops, which is physically impossible.

As stated previous, this is a prototype and limitation as this are not necessarily concerns
for games development. Keep in mind that this algorithm is synchronous, but it's totally
possible to paralelize it.
*/

func (ba *HumidityBoard) Draw(screen *ebiten.Image) {
	var clr color.Color
	for x, row := range ba.values {
		for y, v0 := range row {
			if ba.Rocks[x][y] {
				clr = color.RGBA{255, 255, 255, uint8((v0[1] / math.MaxFloat32) * 255)}
			} else {

				clr = util.HSVColor{uint16((v0[0] / 1023) * 240), 255, 255}
			}
			if x == ba.hvrX && y == ba.hvrY {
				clr = color.Black
			}
			screen.Set(x, y, clr)
		}
	}
}

func (ba *HumidityBoard) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 0, 0
}

func (ba *HumidityBoard) Setup(init [][]mgl32.Vec2) {
	ba.initValues = init
	ba.values = init
}

func (ba *HumidityBoard) Reset() error {
	ba.values = ba.initValues
	return nil
}

func (ba *HumidityBoard) GetState() [][]mgl32.Vec2 {
	return ba.values
}

func (ba *HumidityBoard) Click(btn ebiten.MouseButton) {
	if btn == ebiten.MouseButtonLeft {
		ba.Rocks[ba.hvrX][ba.hvrY] = true
	} else if btn == ebiten.MouseButtonRight {
		ba.Rocks[ba.hvrX][ba.hvrY] = false
	}
}

func (ba *HumidityBoard) Hover(x, y int) {
	ba.hvrX = x
	ba.hvrY = y
}

func MakeHumidityGrid(rockMask, rainMask [][]bool) [][]mgl32.Vec2 {
	// Generate Grid

	grid := util.MakeMatrixWH(len(rockMask), len(rockMask[0]), mgl32.Vec2{0, 1})

	// Generate Rocks
	util.ApplyMaskOnMatrix(grid, rockMask, mgl32.Vec2{0, math.MaxFloat32})

	// Generate Rain
	util.ApplyMaskOnMatrix(grid, rainMask, mgl32.Vec2{1023, 1})
	return grid
}
